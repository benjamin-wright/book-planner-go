package manager

import (
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/crds"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/resources"
)

type state struct {
	cdbs        bucket[crds.CockroachDB, *crds.CockroachDB]
	cclients    bucket[crds.CockroachClient, *crds.CockroachClient]
	cmigrations bucket[crds.CockroachMigration, *crds.CockroachMigration]
	rdbs        bucket[crds.RedisDB, *crds.RedisDB]
	csss        bucket[resources.CockroachStatefulSet, *resources.CockroachStatefulSet]
}

func newState() state {
	return state{
		cdbs:        newBucket[crds.CockroachDB](),
		cclients:    newBucket[crds.CockroachClient](),
		cmigrations: newBucket[crds.CockroachMigration](),
		rdbs:        newBucket[crds.RedisDB](),
		csss:        newBucket[resources.CockroachStatefulSet](),
	}
}
