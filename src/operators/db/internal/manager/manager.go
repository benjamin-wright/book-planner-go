package manager

import (
	"context"
	"fmt"

	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s"
)

type K8sClient[T any] interface {
	Watch(ctx context.Context, cancel context.CancelFunc) (<-chan map[string]T, error)
}

func Manage(
	cdbClient K8sClient[k8s.CockroachDB],
	ccClient K8sClient[k8s.CockroachClient],
	cmClient K8sClient[k8s.CockroachMigration],
	rdbClient K8sClient[k8s.RedisDB],
) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())

	cdbs, err := cdbClient.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach dbs: %+v", err)
	}

	cclients, err := ccClient.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach clients: %+v", err)
	}

	cmigrations, err := cmClient.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch cockroach migration: %+v", err)
	}

	rdbs, err := rdbClient.Watch(ctx, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to watch redis dbs: %+v", err)
	}

	return cancel, nil
}
