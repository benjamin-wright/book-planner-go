package manager

import (
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/cockroach"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/crds"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/resources"
)

type state struct {
	cdbs         bucket[crds.CockroachDB, *crds.CockroachDB]
	cclients     bucket[crds.CockroachClient, *crds.CockroachClient]
	cmigrations  bucket[crds.CockroachMigration, *crds.CockroachMigration]
	csss         bucket[resources.CockroachStatefulSet, *resources.CockroachStatefulSet]
	cpvcs        bucket[resources.CockroachPVC, *resources.CockroachPVC]
	csvcs        bucket[resources.CockroachService, *resources.CockroachService]
	csecrets     bucket[resources.CockroachSecret, *resources.CockroachSecret]
	cdatabases   bucket[cockroach.Database, *cockroach.Database]
	cusers       bucket[cockroach.User, *cockroach.User]
	cpermissions bucket[cockroach.Permission, *cockroach.Permission]
	rdbs         bucket[crds.RedisDB, *crds.RedisDB]
	rsss         bucket[resources.RedisStatefulSet, *resources.RedisStatefulSet]
	rpvcs        bucket[resources.RedisPVC, *resources.RedisPVC]
	rsvcs        bucket[resources.RedisService, *resources.RedisService]
}

func newState() state {
	return state{
		cdbs:         newBucket[crds.CockroachDB](),
		cclients:     newBucket[crds.CockroachClient](),
		cmigrations:  newBucket[crds.CockroachMigration](),
		csss:         newBucket[resources.CockroachStatefulSet](),
		cpvcs:        newBucket[resources.CockroachPVC](),
		csvcs:        newBucket[resources.CockroachService](),
		csecrets:     newBucket[resources.CockroachSecret](),
		cdatabases:   newBucket[cockroach.Database](),
		cusers:       newBucket[cockroach.User](),
		cpermissions: newBucket[cockroach.Permission](),
		rdbs:         newBucket[crds.RedisDB](),
		rsss:         newBucket[resources.RedisStatefulSet](),
		rpvcs:        newBucket[resources.RedisPVC](),
		rsvcs:        newBucket[resources.RedisService](),
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

func clientDemand[T comparable, PT Nameable[T]](
	s *state,
	existing map[string]T,
	transform func(crds.CockroachClient) T,
) demand[T] {
	d := demand[T]{
		toAdd:    []T{},
		toRemove: []T{},
	}

	seen := map[string]T{}

	for _, client := range s.cclients.state {
		ss, hasSS := s.csss.state[client.Deployment]
		_, hasSvc := s.csvcs.state[client.Deployment]

		if !hasSS || !hasSvc || !ss.Ready {
			continue
		}

		desired := transform(client)
		ptr := PT(&desired)
		seen[ptr.GetName()] = desired

		if _, ok := existing[ptr.GetName()]; !ok {
			d.toAdd = append(d.toAdd, desired)
		}
	}

	for current, db := range existing {
		if _, ok := seen[current]; !ok {
			d.toRemove = append(d.toRemove, db)
		}
	}

	return d
}

func (s *state) getCDBDemand() demand[cockroach.Database] {
	return clientDemand(
		s,
		s.cdatabases.state,
		func(client crds.CockroachClient) cockroach.Database {
			return cockroach.Database{
				Name: client.Database,
				DB:   client.Deployment,
			}
		},
	)
}

func (s *state) getCUserDemand() demand[cockroach.User] {
	return clientDemand(
		s,
		s.cusers.state,
		func(client crds.CockroachClient) cockroach.User {
			return cockroach.User{
				Name: client.Username,
				DB:   client.Deployment,
			}
		},
	)
}

func (s *state) getCPermissionDemand() demand[cockroach.Permission] {
	return clientDemand(
		s,
		s.cpermissions.state,
		func(client crds.CockroachClient) cockroach.Permission {
			return cockroach.Permission{
				User:     client.Username,
				Database: client.Database,
				DB:       client.Deployment,
			}
		},
	)
}

func (s *state) getCSecretsDemand() demand[resources.CockroachSecret] {
	return clientDemand(
		s,
		s.csecrets.state,
		func(client crds.CockroachClient) resources.CockroachSecret {
			return resources.CockroachSecret{
				Name:     client.Secret,
				DB:       client.Deployment,
				Database: client.Database,
				User:     client.Username,
			}
		},
	)
}

func (s *state) getCMigrationsDemand() map[string]map[string]map[int64]crds.CockroachMigration {
	migrations := map[string]map[string]map[int64]crds.CockroachMigration{}

	// Get migrations as mapped lookup
	for _, migration := range s.cmigrations.state {
		if _, ok := migrations[migration.Deployment]; !ok {
			migrations[migration.Deployment] = map[string]map[int64]crds.CockroachMigration{}
		}

		if _, ok := migrations[migration.Deployment][migration.Database]; !ok {
			migrations[migration.Deployment][migration.Database] = map[int64]crds.CockroachMigration{}
		}

		migrations[migration.Deployment][migration.Database][migration.Index] = migration
	}

	// Pick out the migrations for statefulsets that
	demand := map[string]map[string]map[int64]crds.CockroachMigration{}
	for _, db := range s.cdatabases.state {
		if dbMigrations, ok := migrations[db.DB]; ok {
			if migrations, ok := dbMigrations[db.Name]; ok {
				if _, ok := demand[db.DB]; !ok {
					demand[db.DB] = map[string]map[int64]crds.CockroachMigration{}
				}

				demand[db.DB][db.Name] = migrations
			}
		}
	}

	return demand
}
