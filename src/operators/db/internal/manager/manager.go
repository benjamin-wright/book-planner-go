package manager

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/manager/state"
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
	csss        K8sClient[resources.CockroachStatefulSet]
	cpvcs       K8sClient[resources.CockroachPVC]
	csvcs       K8sClient[resources.CockroachService]
	csecrets    K8sClient[resources.CockroachSecret]
	rdbs        K8sClient[crds.RedisDB]
	rsss        K8sClient[resources.RedisStatefulSet]
	rpvcs       K8sClient[resources.RedisPVC]
	rsvcs       K8sClient[resources.RedisService]
}

type streams struct {
	cdbs        <-chan k8s_generic.Update[crds.CockroachDB]
	cclients    <-chan k8s_generic.Update[crds.CockroachClient]
	cmigrations <-chan k8s_generic.Update[crds.CockroachMigration]
	csss        <-chan k8s_generic.Update[resources.CockroachStatefulSet]
	cpvcs       <-chan k8s_generic.Update[resources.CockroachPVC]
	csvcs       <-chan k8s_generic.Update[resources.CockroachService]
	csecrets    <-chan k8s_generic.Update[resources.CockroachSecret]
	rdbs        <-chan k8s_generic.Update[crds.RedisDB]
	rsss        <-chan k8s_generic.Update[resources.RedisStatefulSet]
	rpvcs       <-chan k8s_generic.Update[resources.RedisPVC]
	rsvcs       <-chan k8s_generic.Update[resources.RedisService]
}

type Manager struct {
	namespace string
	ctx       context.Context
	cancel    context.CancelFunc
	clients   clients
	streams   streams
	state     state.State
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
	cssClient K8sClient[resources.CockroachStatefulSet],
	cpvcClient K8sClient[resources.CockroachPVC],
	csvcClient K8sClient[resources.CockroachService],
	csecretClient K8sClient[resources.CockroachSecret],
	rdbClient K8sClient[crds.RedisDB],
	rssClient K8sClient[resources.RedisStatefulSet],
	rpvcClient K8sClient[resources.RedisPVC],
	rsvcClient K8sClient[resources.RedisService],
	debouncer time.Duration,
) (*Manager, error) {
	ctx, cancel := context.WithCancel(context.Background())

	clients := clients{
		cdbs:        cdbClient,
		cclients:    ccClient,
		cmigrations: cmClient,
		csss:        cssClient,
		cpvcs:       cpvcClient,
		csvcs:       csvcClient,
		csecrets:    csecretClient,
		rdbs:        rdbClient,
		rsss:        rssClient,
		rpvcs:       rpvcClient,
		rsvcs:       rsvcClient,
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

	csecrets, err := clients.csecrets.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach secrets: %+v", err)
	}

	rdbs, err := clients.rdbs.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch redis dbs: %+v", err)
	}

	rsss, err := clients.rsss.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch redis stateful sets: %+v", err)
	}

	rpvcs, err := clients.rpvcs.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch redis persistent volume claims: %+v", err)
	}

	rsvcs, err := clients.rsvcs.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch redis services: %+v", err)
	}

	streams := streams{
		cdbs:        cdbs,
		cclients:    cclients,
		cmigrations: cmigrations,
		csss:        csss,
		cpvcs:       cpvcs,
		csvcs:       csvcs,
		csecrets:    csecrets,
		rdbs:        rdbs,
		rsss:        rsss,
		rpvcs:       rpvcs,
		rsvcs:       rsvcs,
	}

	return &Manager{
		namespace: namespace,
		ctx:       ctx,
		cancel:    cancel,
		clients:   clients,
		streams:   streams,
		state:     state.New(),
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
		m.state.Apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.cclients:
		m.state.Apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.cmigrations:
		m.state.Apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.rdbs:
		m.state.Apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.csss:
		m.state.Apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.csvcs:
		m.state.Apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.cpvcs:
		m.state.Apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.csecrets:
		m.state.Apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.rdbs:
		m.state.Apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.rsss:
		m.state.Apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.rpvcs:
		m.state.Apply(update)
		m.debouncer.Trigger()
	case update := <-m.streams.rsvcs:
		m.state.Apply(update)
		m.debouncer.Trigger()
	case <-m.debouncer.Wait():
		zap.S().Infof("Database state update")
		m.state.RefreshCockroach(m.namespace)
		zap.S().Infof("Processing Started")
		m.processCockroachDBs()
		m.processCockroachClients()
		m.processCockroachMigrations()
		m.processRedisDBs()
		zap.S().Infof("Processing Done")
	}
}

func (m *Manager) processCockroachDBs() {
	ssDemand := m.state.GetCSSSDemand()
	svcDemand := m.state.GetCSvcDemand()
	pvcsToRemove := m.state.GetCPVCDemand()

	for _, db := range ssDemand.ToRemove {
		zap.S().Infof("Deleting db: %s", db.Name)
		err := m.clients.csss.Delete(m.ctx, db.Name)

		if err != nil {
			zap.S().Errorf("Failed to delete cockroachdb stateful set: %+v", err)
		}
	}

	for _, svc := range svcDemand.ToRemove {
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

	for _, db := range ssDemand.ToAdd {
		zap.S().Infof("Creating db: %s", db.Name)
		err := m.clients.csss.Create(m.ctx, db)

		if err != nil {
			zap.S().Errorf("Failed to create cockroachdb stateful set: %+v", err)
		}
	}

	for _, svc := range svcDemand.ToAdd {
		zap.S().Infof("Creating service: %s", svc.Name)
		err := m.clients.csvcs.Create(m.ctx, svc)

		if err != nil {
			zap.S().Errorf("Failed to create cockroachdb service: %+v", err)
		}
	}
}

func (m *Manager) processCockroachClients() {
	dbDemand := m.state.GetCDBDemand()
	userDemand := m.state.GetCUserDemand()
	permsDemand := m.state.GetCPermissionDemand()
	secretsDemand := m.state.GetCSecretsDemand()

	dbs := map[string]struct{}{}
	for _, db := range dbDemand.ToAdd {
		dbs[db.DB] = struct{}{}
	}
	for _, db := range dbDemand.ToRemove {
		dbs[db.DB] = struct{}{}
	}
	for _, user := range userDemand.ToAdd {
		dbs[user.DB] = struct{}{}
	}
	for _, user := range userDemand.ToRemove {
		dbs[user.DB] = struct{}{}
	}
	for _, perm := range permsDemand.ToAdd {
		dbs[perm.DB] = struct{}{}
	}
	for _, perm := range permsDemand.ToRemove {
		dbs[perm.DB] = struct{}{}
	}

	for _, secret := range secretsDemand.ToRemove {
		zap.S().Infof("Removing secret %s", secret.Name)
		err := m.clients.csecrets.Delete(m.ctx, secret.Name)
		if err != nil {
			zap.S().Errorf("Failed to delete secret %s: %+v", secret.Name, err)
		}
	}

	for database := range dbs {
		cli, err := cockroach.New(database, m.namespace)
		if err != nil {
			zap.S().Errorf("Failed to create database client for %s: %+v", database, err)
			continue
		}
		defer cli.Stop()

		for _, perm := range permsDemand.ToRemove {
			if perm.DB != database {
				continue
			}

			zap.S().Infof("Dropping permission for user %s in database %s of db %s", perm.User, perm.Database, perm.DB)
			err = cli.RevokePermission(perm)
			if err != nil {
				zap.S().Errorf("Failed to revoke permission: %+v", err)
			}
		}

		for _, db := range dbDemand.ToRemove {
			if db.DB != database {
				continue
			}

			zap.S().Infof("Dropping database %s in db %s", db.Name, db.DB)
			err = cli.DeleteDB(db)
			if err != nil {
				zap.S().Errorf("Failed to delete database: %+v", err)
			}
		}

		for _, user := range userDemand.ToRemove {
			if user.DB != database {
				continue
			}

			zap.S().Infof("Dropping user %s in db %s", user.Name, user.DB)
			err = cli.DeleteUser(user)
			if err != nil {
				zap.S().Errorf("Failed to delete user: %+v", err)
			}
		}

		for _, db := range dbDemand.ToAdd {
			if db.DB != database {
				continue
			}

			zap.S().Infof("Creating database %s in db %s", db.Name, db.DB)

			err := cli.CreateDB(db)
			if err != nil {
				zap.S().Errorf("Failed to create database: %+v", err)
			}
		}

		for _, user := range userDemand.ToAdd {
			if user.DB != database {
				continue
			}

			zap.S().Infof("Creating user %s in db %s", user.Name, user.DB)

			err := cli.CreateUser(user)
			if err != nil {
				zap.S().Errorf("Failed to create user: %+v", err)
			}
		}

		for _, perm := range permsDemand.ToAdd {
			if perm.DB != database {
				continue
			}

			zap.S().Infof("Adding permission for user %s in database %s of db %s", perm.User, perm.Database, perm.DB)
			err := cli.GrantPermission(perm)
			if err != nil {
				zap.S().Errorf("Failed to grant permission: %+v", err)
			}
		}
	}

	for _, secret := range secretsDemand.ToAdd {
		zap.S().Infof("Adding secret %s", secret.Name)
		err := m.clients.csecrets.Create(m.ctx, secret)
		if err != nil {
			zap.S().Errorf("Failed to create secret %s: %+v", secret.Name, err)
		}
	}
}

func (m *Manager) processCockroachMigrations() {
	demand := m.state.GetCMigrationsDemand()

	for _, deployment := range demand.GetDBs() {
		for _, database := range demand.GetDatabases(deployment) {
			client, err := cockroach.NewMigrations(deployment, m.namespace, database)
			if err != nil {
				zap.S().Errorf("Failed to create migrations client: %+v", err)
				continue
			}
			defer client.Stop()

			for demand.Next(deployment, database) {
				migration, index := demand.GetNextMigration(deployment, database)

				zap.S().Infof("Running migration %s [%s] %d", deployment, database, index)

				err := client.RunMigration(index, migration)
				if err != nil {
					zap.S().Errorf("Failed to run migration %d: %+v", index, err)
					break
				}
			}
		}
	}
}

func (m *Manager) processRedisDBs() {
	ssDemand := m.state.GetRSSSDemand()
	svcDemand := m.state.GetRSvcDemand()
	pvcsToRemove := m.state.GetRPVCDemand()

	for _, db := range ssDemand.ToRemove {
		zap.S().Infof("Deleting db: %s", db.Name)
		err := m.clients.rsss.Delete(m.ctx, db.Name)

		if err != nil {
			zap.S().Errorf("Failed to delete redis stateful set: %+v", err)
		}
	}

	for _, svc := range svcDemand.ToRemove {
		zap.S().Infof("Deleting service: %s", svc.Name)
		err := m.clients.rsvcs.Delete(m.ctx, svc.Name)

		if err != nil {
			zap.S().Errorf("Failed to delete redis service: %+v", err)
		}
	}

	for _, pvc := range pvcsToRemove {
		zap.S().Infof("Deleting pvc: %s", pvc.Name)
		err := m.clients.rpvcs.Delete(m.ctx, pvc.Name)

		if err != nil {
			zap.S().Errorf("Failed to delete redis PVC: %+v", err)
		}
	}

	for _, db := range ssDemand.ToAdd {
		zap.S().Infof("Creating db: %s", db.Name)
		err := m.clients.rsss.Create(m.ctx, db)

		if err != nil {
			zap.S().Errorf("Failed to create redis stateful set: %+v", err)
		}
	}

	for _, svc := range svcDemand.ToAdd {
		zap.S().Infof("Creating service: %s", svc.Name)
		err := m.clients.rsvcs.Create(m.ctx, svc)

		if err != nil {
			zap.S().Errorf("Failed to create redis service: %+v", err)
		}
	}
}
