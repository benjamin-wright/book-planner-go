package manager

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s"
)

type demandState struct {
	cdbs        map[string]k8s.CockroachDB
	cclients    map[string]k8s.CockroachClient
	cmigrations map[string]k8s.CockroachMigration
	rdbs        map[string]k8s.RedisDB
}

func newDemandState() demandState {
	return demandState{
		cdbs:        map[string]k8s.CockroachDB{},
		cclients:    map[string]k8s.CockroachClient{},
		cmigrations: map[string]k8s.CockroachMigration{},
		rdbs:        map[string]k8s.RedisDB{},
	}
}

func (s *demandState) log() {
	maps := []string{}
	maps = append(maps, "\nCockroach DBs", logMap(s.cdbs))
	maps = append(maps, "Cockroach Clients", logMap(s.cclients))
	maps = append(maps, "Cockroach Migrations", logMap(s.cmigrations))
	maps = append(maps, "Redis DBs", logMap(s.rdbs))

	zap.S().Infof("State: %s", strings.Join(maps, "\n"))
}

func logMap[T any](data map[string]T) string {
	lines := []string{}

	for key, item := range data {
		lines = append(lines, fmt.Sprintf(" - %s: %+v", key, item))
	}

	return strings.Join(lines, "\n")
}
