package manager

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/crds"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/resources"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

type Manager struct {
	cdbs        K8sClient[crds.CockroachDB]
	cclients    K8sClient[crds.CockroachClient]
	cmigrations K8sClient[crds.CockroachMigration]
	rdbs        K8sClient[crds.RedisDB]
	csss        K8sClient[resources.CockroachStatefulSet]
}

type K8sClient[T any] interface {
	Watch(ctx context.Context, cancel context.CancelFunc) (<-chan k8s_generic.Update[T], error)
}

func New(
	cdbClient K8sClient[crds.CockroachDB],
	ccClient K8sClient[crds.CockroachClient],
	cmClient K8sClient[crds.CockroachMigration],
	rdbClient K8sClient[crds.RedisDB],
	cssClient K8sClient[resources.CockroachStatefulSet],
) *Manager {
	return &Manager{
		cdbs:        cdbClient,
		cclients:    ccClient,
		cmigrations: cmClient,
		rdbs:        rdbClient,
		csss:        cssClient,
	}
}

func (m *Manager) Start() (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())

	cdbs, err := m.cdbs.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach dbs: %+v", err)
	}

	cclients, err := m.cclients.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach clients: %+v", err)
	}

	cmigrations, err := m.cmigrations.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach migration: %+v", err)
	}

	rdbs, err := m.rdbs.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch redis dbs: %+v", err)
	}

	csss, err := m.csss.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach stateful sets: %+v", err)
	}

	state := newState()

	go func() {
		for {
			select {
			case update := <-cdbs:
				state.cdbs.apply(update)
			case update := <-cclients:
				state.cclients.apply(update)
			case update := <-cmigrations:
				state.cmigrations.apply(update)
			case update := <-rdbs:
				state.rdbs.apply(update)
			case update := <-csss:
				state.csss.apply(update)
			case <-ctx.Done():
				zap.S().Infof("context cancelled, exiting manager loop")
				return
			}

			zap.S().Infof("%+v", state)
		}
	}()

	return cancel, nil
}
