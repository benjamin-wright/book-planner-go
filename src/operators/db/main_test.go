//go:build !unit
// +build !unit

package main

import (
	"context"
	"os"
	"testing"

	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/crds"
	"ponglehub.co.uk/book-planner-go/src/pkg/k8s_generic"
)

func makeClients(t *testing.T, namespace string) (
	*k8s_generic.Client[crds.CockroachDB, *crds.CockroachDB],
	*k8s_generic.Client[crds.CockroachClient, *crds.CockroachClient],
) {
	cdbs, err := crds.NewCockroachDBClient(namespace)
	if err != nil {
		t.Logf("failed to create cdb client: %+v", err)
		t.FailNow()
	}

	cclients, err := crds.NewCockroachClientClient(namespace)
	if err != nil {
		t.Logf("failed to create cclient client: %+v", err)
		t.FailNow()
	}

	return cdbs, cclients
}

func TestSomethingIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	namespace := os.Getenv("NAMESPACE")

	cdbs, cclients := makeClients(t, namespace)

	err := cdbs.DeleteAll(context.Background())
	if err != nil {
		t.Logf("failed to clear existing dbs: %+v", err)
		t.FailNow()
	}

	err = cclients.DeleteAll(context.Background())
	if err != nil {
		t.Logf("failed to clear existing clients: %+v", err)
		t.FailNow()
	}

	err = cdbs.Create(context.Background(), crds.CockroachDB{
		Name:    "other-db",
		Storage: "512Mi",
	})
	if err != nil {
		t.Logf("failed to create test db: %+v", err)
		t.FailNow()
	}

	err = cclients.Create(context.Background(), crds.CockroachClient{
		Deployment: "other-db",
		Database:   "my_db",
		Name:       "my-client",
		Username:   "my_user",
		Secret:     "my-secret",
	})
	if err != nil {
		t.Logf("failed to create test client: %+v", err)
		t.FailNow()
	}
}
