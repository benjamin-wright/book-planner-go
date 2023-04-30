//go:build !unit
// +build !unit

package main

import (
	"context"
	"os"
	"testing"

	"ponglehub.co.uk/book-planner-go/src/operators/db/pkg/k8s"
)

func TestSomethingIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	namespace := os.Getenv("NAMESPACE")

	cli, err := k8s.NewCockroachDBClient(namespace)
	if err != nil {
		t.Logf("failed to create kube client: %+v", err)
		t.FailNow()
	}

	err = cli.DeleteAll(context.Background())
	if err != nil {
		t.Logf("failed to clear existing dbs: %+v", err)
		t.FailNow()
	}

	err = cli.Create(context.Background(), k8s.CockroachDB{
		Name:      "test-db",
		Namespace: namespace,
		Storage:   "1Gi",
	})
	if err != nil {
		t.Logf("failed to create test db: %+v", err)
		t.FailNow()
	}

	err = cli.Create(context.Background(), k8s.CockroachDB{
		Name:      "other-db",
		Namespace: namespace,
		Storage:   "512Mi",
	})
	if err != nil {
		t.Logf("failed to create other db: %+v", err)
		t.FailNow()
	}
}
