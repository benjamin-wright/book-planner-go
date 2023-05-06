package manager

import (
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/cockroach"
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
	csvcs       bucket[resources.CockroachService, *resources.CockroachService]
	cdatabases  bucket[cockroach.Database, *cockroach.Database]
}

func newState() state {
	return state{
		cdbs:        newBucket[crds.CockroachDB](),
		cclients:    newBucket[crds.CockroachClient](),
		cmigrations: newBucket[crds.CockroachMigration](),
		rdbs:        newBucket[crds.RedisDB](),
		csss:        newBucket[resources.CockroachStatefulSet](),
		cpvcs:       newBucket[resources.CockroachPVC](),
		csvcs:       newBucket[resources.CockroachService](),
		cdatabases:  newBucket[cockroach.Database](),
	}
}

type demand[T any] struct {
	toAdd    []T
	toRemove []T
}

func (s *state) getCSSSDemand() demand[resources.CockroachStatefulSet] {
	toAdd := []resources.CockroachStatefulSet{}
	toRemove := []resources.CockroachStatefulSet{}

	for name, db := range s.cdbs.state {
		if ss, ok := s.csss.state[name]; !ok {
			toAdd = append(toAdd, resources.CockroachStatefulSet{
				Name:    db.Name,
				Storage: db.Storage,
			})
		} else if db.Storage != ss.Storage {
			toRemove = append(toRemove, resources.CockroachStatefulSet{
				Name: db.Name,
			})
			toAdd = append(toAdd, resources.CockroachStatefulSet{
				Name:    db.Name,
				Storage: db.Storage,
			})
		}
	}

	for name := range s.csss.state {
		if _, ok := s.cdbs.state[name]; !ok {
			toRemove = append(toRemove, resources.CockroachStatefulSet{
				Name: name,
			})
		}
	}

	return demand[resources.CockroachStatefulSet]{
		toAdd:    toAdd,
		toRemove: toRemove,
	}
}

func (s *state) getCSvcDemand() demand[resources.CockroachService] {
	toAdd := []resources.CockroachService{}
	toRemove := []resources.CockroachService{}

	for name := range s.cdbs.state {
		if _, ok := s.csvcs.state[name]; !ok {
			toAdd = append(toAdd, resources.CockroachService{
				Name: name,
			})
		}
	}

	for name := range s.csss.state {
		if _, ok := s.cdbs.state[name]; !ok {
			toRemove = append(toRemove, resources.CockroachService{
				Name: name,
			})
		}
	}

	return demand[resources.CockroachService]{
		toAdd:    toAdd,
		toRemove: toRemove,
	}
}

func (s *state) getCPVCDemand(toRemove []resources.CockroachStatefulSet) []resources.CockroachPVC {
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

func (s *state) getCDBDemand() demand[cockroach.Database] {
	d := demand[cockroach.Database]{
		toAdd:    []cockroach.Database{},
		toRemove: []cockroach.Database{},
	}

	seen := map[string]cockroach.Database{}

	for _, client := range s.cclients.state {
		ss, hasSS := s.csss.state[client.Deployment]
		_, hasSvc := s.csvcs.state[client.Deployment]

		if !hasSS || !hasSvc || !ss.Ready {
			continue
		}

		desired := cockroach.Database{
			Name: client.Database,
			DB:   client.Deployment,
		}
		seen[desired.GetName()] = desired

		if _, ok := s.cdatabases.state[desired.GetName()]; !ok {
			d.toAdd = append(d.toAdd, desired)
		}
	}

	for current, db := range s.cdatabases.state {
		if _, ok := seen[current]; !ok {
			d.toRemove = append(d.toRemove, db)
		}
	}

	return d
}
