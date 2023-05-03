package manager

import (
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s"
)

type demandState struct {
	cdbs        bucket[k8s.CockroachDB, *k8s.CockroachDB]
	cclients    bucket[k8s.CockroachClient, *k8s.CockroachClient]
	cmigrations bucket[k8s.CockroachMigration, *k8s.CockroachMigration]
	rdbs        bucket[k8s.RedisDB, *k8s.RedisDB]
}

func newDemandState() demandState {
	return demandState{
		cdbs:        newBucket[k8s.CockroachDB](),
		cclients:    newBucket[k8s.CockroachClient](),
		cmigrations: newBucket[k8s.CockroachMigration](),
		rdbs:        newBucket[k8s.RedisDB](),
	}
}
