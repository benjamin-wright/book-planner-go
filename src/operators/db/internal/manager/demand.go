package manager

import (
	"context"

	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/k8s"
)

func Manage(
	ctx context.Context,
	cdbs <-chan map[string]k8s.CockroachDB,
	cclients <-chan map[string]k8s.CockroachClient,
	cmigrations <-chan map[string]k8s.CockroachMigration,
	rdbs <-chan map[string]k8s.RedisDB,
) {
	state := newDemandState()

	go func(
		done <-chan struct{},
		cdbs <-chan map[string]k8s.CockroachDB,
		cclients <-chan map[string]k8s.CockroachClient,
		cmigrations <-chan map[string]k8s.CockroachMigration,
		rdbs <-chan map[string]k8s.RedisDB,
	) {
		for {
			select {
			case dbs := <-cdbs:
				state.cdbs = dbs
			case clients := <-cclients:
				state.cclients = clients
			case migrations := <-cmigrations:
				state.cmigrations = migrations
			case dbs := <-rdbs:
				state.rdbs = dbs
			case <-done:
				return
			}

			state.log()
		}
	}(ctx.Done(), cdbs, cclients, cmigrations, rdbs)
}
