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

func (s *state) getCSSSDemand() demand[resources.CockroachStatefulSet] {
	return getStorageBoundDemand(
		s.cdbs.state,
		s.csss.state,
		func(db crds.CockroachDB) resources.CockroachStatefulSet {
			return resources.CockroachStatefulSet{
				Name:    db.Name,
				Storage: db.Storage,
			}
		},
	)
}

func (s *state) getCSvcDemand() demand[resources.CockroachService] {
	return getOneForOneDemand(
		s.cdbs.state,
		s.csvcs.state,
		func(db crds.CockroachDB) resources.CockroachService {
			return resources.CockroachService{Name: db.Name}
		},
	)
}

func (s *state) getCPVCDemand() []resources.CockroachPVC {
	return getOrphanedDemand(
		s.csss.state,
		s.cpvcs.state,
		func(ss resources.CockroachStatefulSet, pvc resources.CockroachPVC) bool {
			return ss.Name == pvc.Database
		},
	)
}

func (s *state) getCDBDemand() demand[cockroach.Database] {
	return getServiceBoundDemand(
		s.cclients.state,
		s.cdatabases.state,
		s.csss.state,
		s.csvcs.state,
		func(client crds.CockroachClient) cockroach.Database {
			return cockroach.Database{
				Name: client.Database,
				DB:   client.Deployment,
			}
		},
	)
}

func (s *state) getCUserDemand() demand[cockroach.User] {
	return getServiceBoundDemand(
		s.cclients.state,
		s.cusers.state,
		s.csss.state,
		s.csvcs.state,
		func(client crds.CockroachClient) cockroach.User {
			return cockroach.User{
				Name: client.Username,
				DB:   client.Deployment,
			}
		},
	)
}

func (s *state) getCPermissionDemand() demand[cockroach.Permission] {
	return getServiceBoundDemand(
		s.cclients.state,
		s.cpermissions.state,
		s.csss.state,
		s.csvcs.state,
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
	return getServiceBoundDemand(
		s.cclients.state,
		s.csecrets.state,
		s.csss.state,
		s.csvcs.state,
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

func (s *state) getRSSSDemand() demand[resources.RedisStatefulSet] {
	return getStorageBoundDemand(
		s.rdbs.state,
		s.rsss.state,
		func(db crds.RedisDB) resources.RedisStatefulSet {
			return resources.RedisStatefulSet{
				Name:    db.Name,
				Storage: db.Storage,
			}
		},
	)
}

func (s *state) getRSvcDemand() demand[resources.RedisService] {
	return getOneForOneDemand(
		s.rdbs.state,
		s.rsvcs.state,
		func(db crds.RedisDB) resources.RedisService {
			return resources.RedisService{Name: db.Name}
		},
	)
}

func (s *state) getRPVCDemand() []resources.RedisPVC {
	return getOrphanedDemand(
		s.rsss.state,
		s.rpvcs.state,
		func(ss resources.RedisStatefulSet, pvc resources.RedisPVC) bool {
			return ss.Name == pvc.Database
		},
	)
}
