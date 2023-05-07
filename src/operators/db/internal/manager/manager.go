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
	csecrets    K8sClient[resources.CockroachSecret]
}

type streams struct {
	cdbs        <-chan k8s_generic.Update[crds.CockroachDB]
	cclients    <-chan k8s_generic.Update[crds.CockroachClient]
	cmigrations <-chan k8s_generic.Update[crds.CockroachMigration]
	rdbs        <-chan k8s_generic.Update[crds.RedisDB]
	csss        <-chan k8s_generic.Update[resources.CockroachStatefulSet]
	cpvcs       <-chan k8s_generic.Update[resources.CockroachPVC]
	csvcs       <-chan k8s_generic.Update[resources.CockroachService]
	csecrets    <-chan k8s_generic.Update[resources.CockroachSecret]
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
	csecretClient K8sClient[resources.CockroachSecret],
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
		csecrets:    csecretClient,
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

	csecrets, err := clients.csecrets.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach secrets: %+v", err)
	}

	streams := streams{
		cdbs:        cdbs,
		cclients:    cclients,
		cmigrations: cmigrations,
		rdbs:        rdbs,
		csss:        csss,
		cpvcs:       cpvcs,
		csvcs:       csvcs,
		csecrets:    csecrets,
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
	case update := <-m.streams.csecrets:
		m.state.csecrets.apply(update)
		m.debouncer.Trigger()
	case <-m.debouncer.Wait():
		zap.S().Infof("Processing Started")
		m.refreshDBs()
		m.processCockroachDBs()
		m.processCockroachClients()
		m.processCockroachMigrations()
		zap.S().Infof("Processing Done")
	}
}

func (m *Manager) refreshDBs() {
	m.state.cdatabases.clear()
	m.state.cusers.clear()
	m.state.cpermissions.clear()

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
		defer cli.Stop()

		users, err := cli.ListUsers()
		if err != nil {
			zap.S().Errorf("Failed to list users in %s: %+v", db, err)
		}

		for _, user := range users {
			m.state.cusers.add(user)
		}

		names, err := cli.ListDBs()
		if err != nil {
			zap.S().Errorf("Failed to list databases in %s: %+v", db, err)
		}

		for _, db := range names {
			m.state.cdatabases.add(db)

			permissions, err := cli.ListPermitted(db)
			if err != nil {
				zap.S().Errorf("Failed to list permissions in %s: %+v", db.Name, err)
			}

			for _, p := range permissions {
				m.state.cpermissions.add(p)
			}
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
	dbDemand := m.state.getCDBDemand()
	userDemand := m.state.getCUserDemand()
	permsDemand := m.state.getCPermissionDemand()
	secretsDemand := m.state.getCSecretsDemand()

	dbs := map[string]struct{}{}
	for _, db := range dbDemand.toAdd {
		dbs[db.DB] = struct{}{}
	}
	for _, db := range dbDemand.toRemove {
		dbs[db.DB] = struct{}{}
	}
	for _, user := range userDemand.toAdd {
		dbs[user.DB] = struct{}{}
	}
	for _, user := range userDemand.toRemove {
		dbs[user.DB] = struct{}{}
	}
	for _, perm := range permsDemand.toAdd {
		dbs[perm.DB] = struct{}{}
	}
	for _, perm := range permsDemand.toRemove {
		dbs[perm.DB] = struct{}{}
	}

	for _, secret := range secretsDemand.toRemove {
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

		for _, perm := range permsDemand.toRemove {
			if perm.DB != database {
				continue
			}

			zap.S().Infof("Dropping permission for user %s in database %s of db %s", perm.User, perm.Database, perm.DB)
			err = cli.RevokePermission(perm)
			if err != nil {
				zap.S().Errorf("Failed to revoke permission: %+v", err)
			}
		}

		for _, db := range dbDemand.toRemove {
			if db.DB != database {
				continue
			}

			zap.S().Infof("Dropping database %s in db %s", db.Name, db.DB)
			err = cli.DeleteDB(db)
			if err != nil {
				zap.S().Errorf("Failed to delete database: %+v", err)
			}
		}

		for _, user := range userDemand.toRemove {
			if user.DB != database {
				continue
			}

			zap.S().Infof("Dropping user %s in db %s", user.Name, user.DB)
			err = cli.DeleteUser(user)
			if err != nil {
				zap.S().Errorf("Failed to delete user: %+v", err)
			}
		}

		for _, db := range dbDemand.toAdd {
			if db.DB != database {
				continue
			}

			zap.S().Infof("Creating database %s in db %s", db.Name, db.DB)

			err := cli.CreateDB(db)
			if err != nil {
				zap.S().Errorf("Failed to create database: %+v", err)
			}
		}

		for _, user := range userDemand.toAdd {
			if user.DB != database {
				continue
			}

			zap.S().Infof("Creating user %s in db %s", user.Name, user.DB)

			err := cli.CreateUser(user)
			if err != nil {
				zap.S().Errorf("Failed to create user: %+v", err)
			}
		}

		for _, perm := range permsDemand.toAdd {
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

	for _, secret := range secretsDemand.toAdd {
		zap.S().Infof("Adding secret %s", secret.Name)
		err := m.clients.csecrets.Create(m.ctx, secret)
		if err != nil {
			zap.S().Errorf("Failed to create secret %s: %+v", secret.Name, err)
		}
	}
}

func (m *Manager) processCockroachMigrations() {
	migrationDemand := m.state.getCMigrationsDemand()

	for deployment, byDatabase := range migrationDemand {
		for database, migrations := range byDatabase {
			client, err := cockroach.NewMigrations(deployment, m.namespace, database)
			if err != nil {
				zap.S().Errorf("Failed to create migrations client: %+v", err)
				continue
			}
			defer client.Stop()

			if ok, err := client.HasMigrationsTable(); err != nil {
				zap.S().Errorf("Failed to check for migrations table: %+v", err)
				continue
			} else if !ok {
				err := client.CreateMigrationsTable()
				if err != nil {
					zap.S().Errorf("Failed to create migrations client: %+v", err)
					continue
				}
			}

			nextIndex := client.LatestMigration() + 1
			for {
				migration, ok := migrations[nextIndex]
				if !ok {
					break
				}

				zap.S().Infof("Running migration %s [%s] %d", migration.Deployment, migration.Database, nextIndex)

				err := client.RunMigration(migration.Migration)
				if err != nil {
					zap.S().Errorf("Failed to run migration %d: %+v", nextIndex, err)
					break
				}

				err = client.AddMigration(nextIndex)
				if err != nil {
					zap.S().Errorf("Failed to add migration index %d: %+v", nextIndex, err)
					break
				}

				nextIndex += 1
			}
		}
	}
}
