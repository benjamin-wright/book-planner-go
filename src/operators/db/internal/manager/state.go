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
	cpvcs       bucket[resources.CockroachPVC, *resources.CockroachPVC]
}

func newState() state {
	return state{
		cdbs:        newBucket[crds.CockroachDB](),
		cclients:    newBucket[crds.CockroachClient](),
		cmigrations: newBucket[crds.CockroachMigration](),
		rdbs:        newBucket[crds.RedisDB](),
		csss:        newBucket[resources.CockroachStatefulSet](),
		cpvcs:       newBucket[resources.CockroachPVC](),
	}
}

type demand[T any] struct {
	toAdd    []T
	toRemove []T
}

func (s *state) getCSSSDemand() demand[crds.CockroachDB] {
	toAdd := []crds.CockroachDB{}
	toRemove := []crds.CockroachDB{}

	for name, db := range s.cdbs.state {
		if ss, ok := s.csss.state[name]; !ok {
			toAdd = append(toAdd, db)
		} else if db.Storage != ss.Storage {
			toRemove = append(toRemove, db)
			toAdd = append(toAdd, db)
		}
	}

	for name := range s.csss.state {
		if _, ok := s.cdbs.state[name]; !ok {
			toRemove = append(toRemove, crds.CockroachDB{
				Name: name,
			})
		}
	}

	return demand[crds.CockroachDB]{
		toAdd:    toAdd,
		toRemove: toRemove,
	}
}

func (s *state) getCPVCDemand(toRemove []crds.CockroachDB) []resources.CockroachPVC {
	pvcsToRemove := []resources.CockroachPVC{}

	for _, db := range toRemove {
		for _, pvc := range s.cpvcs.state {
			if pvc.Database == db.Name {
				pvcsToRemove = append(pvcsToRemove, pvc)
			}
		}
	}

	return pvcsToRemove
}
