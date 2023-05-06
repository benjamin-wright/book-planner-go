package manager

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/cockroach"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/crds"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/resources"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/utils"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type clients struct {
	cdbs        K8sClient[crds.CockroachDB]
	cclients    K8sClient[crds.CockroachClient]
	cmigrations K8sClient[crds.CockroachMigration]
	rdbs        K8sClient[crds.RedisDB]
	csss        K8sClient[resources.CockroachStatefulSet]
	cpvcs       K8sClient[resources.CockroachPVC]
	csvcs       K8sClient[resources.CockroachService]
}

type streams struct {
	cdbs        <-chan k8s_generic.Update[crds.CockroachDB]
	cclients    <-chan k8s_generic.Update[crds.CockroachClient]
	cmigrations <-chan k8s_generic.Update[crds.CockroachMigration]
	rdbs        <-chan k8s_generic.Update[crds.RedisDB]
	csss        <-chan k8s_generic.Update[resources.CockroachStatefulSet]
	cpvcs       <-chan k8s_generic.Update[resources.CockroachPVC]
	csvcs       <-chan k8s_generic.Update[resources.CockroachService]
}

type Manager struct {
	namespace string
	ctx       context.Context
	cancel    context.CancelFunc
	clients   clients
	streams   streams
	state     state
	debouncer utils.Debouncer
}

type K8sClient[T comparable] interface {
	Watch(ctx context.Context, cancel context.CancelFunc) (<-chan k8s_generic.Update[T], error)
	Create(ctx context.Context, resource T) error
	Delete(ctx context.Context, name string) error
	Update(ctx context.Context, resource T) error
}

type CockroachClient interface {
	CreateDB(cockroach.Database) error
}

func New(
	namespace string,
	cdbClient K8sClient[crds.CockroachDB],
	ccClient K8sClient[crds.CockroachClient],
	cmClient K8sClient[crds.CockroachMigration],
	rdbClient K8sClient[crds.RedisDB],
	cssClient K8sClient[resources.CockroachStatefulSet],
	cpvcClient K8sClient[resources.CockroachPVC],
	csvcClient K8sClient[resources.CockroachService],
	debouncer time.Duration,
) (*Manager, error) {
	ctx, cancel := context.WithCancel(context.Background())

	clients := clients{
		cdbs:        cdbClient,
		cclients:    ccClient,
		cmigrations: cmClient,
		rdbs:        rdbClient,
		csss:        cssClient,
		cpvcs:       cpvcClient,
		csvcs:       csvcClient,
	}

	cdbs, err := clients.cdbs.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach dbs: %+v", err)
	}

	cclients, err := clients.cclients.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach clients: %+v", err)
	}

	cmigrations, err := clients.cmigrations.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach migration: %+v", err)
	}

	rdbs, err := clients.rdbs.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch redis dbs: %+v", err)
	}

	csss, err := clients.csss.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach stateful sets: %+v", err)
	}

	cpvcs, err := clients.cpvcs.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach persistent volume claims: %+v", err)
	}

	csvcs, err := clients.csvcs.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach services: %+v", err)
	}

	streams := streams{
		cdbs:        cdbs,
		cclients:    cclients,
		cmigrations: cmigrations,
		rdbs:        rdbs,
		csss:        csss,
		cpvcs:       cpvcs,
		csvcs:       csvcs,
	}

	return &Manager{
		namespace: namespace,
		ctx:       ctx,
		cancel:    cancel,
		clients:   clients,
		streams:   streams,
		state:     newState(),
		debouncer: utils.NewDebouncer(debouncer),
	}, nil
}

func (m *Manager) Stop() {
	m.cancel()
}

func (m *Manager) Start() {
	go func() {
		for {
			select {
			case <-m.ctx.Done():
				zap.S().Infof("context cancelled, exiting manager loop")
				return
			default:
				m.refresh()
			}
		}
	}()
}

func (m *Manager) refresh() {
	select {
	case <-m.ctx.Done():
	case update := <-m.streams.cdbs:
		m.state.cdbs.apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.cclients:
		m.state.cclients.apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.cmigrations:
		m.state.cmigrations.apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.rdbs:
		m.state.rdbs.apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.csss:
		m.state.csss.apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.csvcs:
		m.state.csvcs.apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.cpvcs:
		m.state.cpvcs.apply(update)
		m.debouncer.Trigger()
	case <-m.debouncer.Wait():
		zap.S().Infof("Processing Started")
		m.refreshDBs()
		m.processCockroachDBs()
		m.processCockroachClients()
		zap.S().Infof("Processing Done")
	}
}

func (m *Manager) refreshDBs() {
	m.state.cdatabases.clear()

	for db := range m.state.cdbs.state {
		ss, hasSS := m.state.csss.state[db]
		_, hasSvc := m.state.csvcs.state[db]
		if !hasSS || !hasSvc || !ss.Ready {
			continue
		}

		cli, err := cockroach.New(db, m.namespace)
		if err != nil {
			zap.S().Errorf("Failed to create client for database %s: %+v", db, err)
		}

		names, err := cli.ListDBs()
		if err != nil {
			zap.S().Errorf("Failed to list databases in %s: %+v", db, err)
		}

		for _, name := range names {
			m.state.cdatabases.add(cockroach.Database{
				Name: name,
				DB:   db,
			})
		}
	}
}

func (m *Manager) processCockroachDBs() {
	ssDemand := m.state.getCSSSDemand()
	svcDemand := m.state.getCSvcDemand()
	pvcsToRemove := m.state.getCPVCDemand(ssDemand.toRemove)

	for _, db := range ssDemand.toRemove {
		zap.S().Infof("Deleting db: %s", db.Name)
		err := m.clients.csss.Delete(m.ctx, db.Name)

		if err != nil {
			zap.S().Errorf("Failed to delete cockroachdb stateful set: %+v", err)
		}
	}

	for _, svc := range svcDemand.toRemove {
		zap.S().Infof("Deleting service: %s", svc.Name)
		err := m.clients.csvcs.Delete(m.ctx, svc.Name)

		if err != nil {
			zap.S().Errorf("Failed to delete cockroachdb service: %+v", err)
		}
	}

	for _, pvc := range pvcsToRemove {
		zap.S().Infof("Deleting pvc: %s", pvc.Name)
		err := m.clients.cpvcs.Delete(m.ctx, pvc.Name)

		if err != nil {
			zap.S().Errorf("Failed to delete cockroachdb PVC: %+v", err)
		}
	}

	for _, db := range ssDemand.toAdd {
		zap.S().Infof("Creating db: %s", db.Name)
		err := m.clients.csss.Create(m.ctx, db)

		if err != nil {
			zap.S().Errorf("Failed to create cockroachdb stateful set: %+v", err)
		}
	}

	for _, svc := range svcDemand.toAdd {
		zap.S().Infof("Creating service: %s", svc.Name)
		err := m.clients.csvcs.Create(m.ctx, svc)

		if err != nil {
			zap.S().Errorf("Failed to create cockroachdb service: %+v", err)
		}
	}
}

func (m *Manager) processCockroachClients() {
	demand := m.state.getCDBDemand()

	dbs := map[string]struct{}{}
	for _, db := range demand.toAdd {
		dbs[db.DB] = struct{}{}
	}
	for _, db := range demand.toRemove {
		dbs[db.DB] = struct{}{}
	}

	for db := range dbs {
		cli, err := cockroach.New(db, m.namespace)
		if err != nil {
			zap.S().Errorf("Failed to create database client for %s: %+v", db, err)
		}

		for _, db := range demand.toAdd {
			zap.S().Infof("Creating database %s in db %s", db.Name, db.DB)

			err := cli.CreateDB(db.Name)
			if err != nil {
				zap.S().Errorf("Failed to create database: %+v", err)
			}
		}

		for _, db := range demand.toRemove {
			zap.S().Infof("Dropping database %s in db %s", db.Name, db.DB)
		}
	}

}
