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

	cli, err := k8s.New()
	if err != nil {
		t.Logf("failed to create kube client: %+v", err)
		t.FailNow()
	}

	err = cli.CockroachDBCreate(context.Background(), k8s.CockroachDB{
		Name:      "test-db",
		Namespace: os.Getenv("NAMESPACE"),
	})
	if err != nil {
		t.Logf("failed to create test db: %+v", err)
		t.FailNow()
	}

	err = cli.CockroachDBCreate(context.Background(), k8s.CockroachDB{
		Name:      "other-db",
		Namespace: os.Getenv("NAMESPACE"),
	})
	if err != nil {
		t.Logf("failed to create other db: %+v", err)
		t.FailNow()
	}
}

func TestSomething(t *testing.T) {
	t.Log("Working!")
}
