package manager

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/crds"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/resources"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/utils"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type Manager struct {
	ctx         context.Context
	cancel      context.CancelFunc
	cdbs        K8sClient[crds.CockroachDB]
	cclients    K8sClient[crds.CockroachClient]
	cmigrations K8sClient[crds.CockroachMigration]
	rdbs        K8sClient[crds.RedisDB]
	csss        K8sClient[resources.CockroachStatefulSet]
}

type K8sClient[T any] interface {
	Watch(ctx context.Context, cancel context.CancelFunc) (<-chan k8s_generic.Update[T], error)
	Create(ctx context.Context, resource T) error
}

func New(
	cdbClient K8sClient[crds.CockroachDB],
	ccClient K8sClient[crds.CockroachClient],
	cmClient K8sClient[crds.CockroachMigration],
	rdbClient K8sClient[crds.RedisDB],
	cssClient K8sClient[resources.CockroachStatefulSet],
) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		ctx:         ctx,
		cancel:      cancel,
		cdbs:        cdbClient,
		cclients:    ccClient,
		cmigrations: cmClient,
		rdbs:        rdbClient,
		csss:        cssClient,
	}
}

func (m *Manager) Stop() {
	m.cancel()
}

func (m *Manager) Start(debounce time.Duration) error {
	cdbs, err := m.cdbs.Watch(m.ctx, m.cancel)
	if err != nil {
		return fmt.Errorf("failed to watch cockroach dbs: %+v", err)
	}

	cclients, err := m.cclients.Watch(m.ctx, m.cancel)
	if err != nil {
		return fmt.Errorf("failed to watch cockroach clients: %+v", err)
	}

	cmigrations, err := m.cmigrations.Watch(m.ctx, m.cancel)
	if err != nil {
		return fmt.Errorf("failed to watch cockroach migration: %+v", err)
	}

	rdbs, err := m.rdbs.Watch(m.ctx, m.cancel)
	if err != nil {
		return fmt.Errorf("failed to watch redis dbs: %+v", err)
	}

	csss, err := m.csss.Watch(m.ctx, m.cancel)
	if err != nil {
		return fmt.Errorf("failed to watch cockroach stateful sets: %+v", err)
	}

	state := newState()
	debouncer := utils.NewDebouncer(debounce)

	go func() {
		for {
			select {
			case update := <-cdbs:
				state.cdbs.apply(update)
				debouncer.Trigger()
				continue
			case update := <-cclients:
				state.cclients.apply(update)
				debouncer.Trigger()
				continue
			case update := <-cmigrations:
				state.cmigrations.apply(update)
				debouncer.Trigger()
				continue
			case update := <-rdbs:
				state.rdbs.apply(update)
				debouncer.Trigger()
				continue
			case update := <-csss:
				state.csss.apply(update)
				debouncer.Trigger()
				continue
			case <-debouncer.Wait():
				m.processCockroachDBs(state)
			case <-m.ctx.Done():
				zap.S().Infof("context cancelled, exiting manager loop")
				return
			}
		}
	}()

	return nil
}

func (m *Manager) processCockroachDBs(state state) {
	for name, db := range state.cdbs.state {
		if _, ok := state.csss.state[name]; !ok {
			err := m.csss.Create(m.ctx, resources.CockroachStatefulSet{
				Name:    name,
				Storage: db.Storage,
				CPU:     "100",
				Memory:  "512Mi",
			})

			if err != nil {
				zap.S().Errorf("Failed to create cockroachdb stateful set: %+v", err)
			}
		}
	}
}
