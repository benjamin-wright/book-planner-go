package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/manager"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/crds"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s/resources"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	namespace := os.Getenv("NAMESPACE")

	zap.S().Info("Starting operator...")

	cdbClient, err := crds.NewCockroachDBClient(namespace)
	if err != nil {
		zap.S().Fatalf("Failed to get client for cockroach dbs: %+v", err)
	}

	ccClient, err := crds.NewCockroachClientClient(namespace)
	if err != nil {
		zap.S().Fatalf("Failed to get client for cockroach clients: %+v", err)
	}

	cmClient, err := crds.NewCockroachMigrationClient(namespace)
	if err != nil {
		zap.S().Fatalf("Failed to get client for cockroach migrations: %+v", err)
	}

	rdbClient, err := crds.NewRedisDBClient(namespace)
	if err != nil {
		zap.S().Fatalf("Failed to get client for redis dbs: %+v", err)
	}

	cssClient, err := resources.NewCockroachStatefulSetClient(namespace)
	if err != nil {
		zap.S().Fatalf("Failed to get client for cockroach stateful sets: %+v", err)
	}

	cpvcClient, err := resources.NewCockroachPVCClient(namespace)
	if err != nil {
		zap.S().Fatalf("Failed to get client for cockroach persistent volume claims: %+v", err)
	}

	csvcClient, err := resources.NewCockroachServiceClient(namespace)
	if err != nil {
		zap.S().Fatalf("Failed to get client for cockroach services: %+v", err)
	}

	csecretClient, err := resources.NewCockroachSecretClient(namespace)
	if err != nil {
		zap.S().Fatalf("Failed to get client for cockroach secrets: %+v", err)
	}

	m, err := manager.New(namespace, cdbClient, ccClient, cmClient, rdbClient, cssClient, cpvcClient, csvcClient, csecretClient, time.Second*5)
	if err != nil {
		zap.S().Fatalf("Failed to start the manager: %+v", err)
	}

	m.Start()
	zap.S().Info("Running!")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")
	m.Stop()
}
