package main

import (
	"log"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/manager"
	"ponglehub.co.uk/book-planner-go/src/operators/db/internal/services/k8s"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	namespace := os.Getenv("NAMESPACE")

	zap.S().Info("Starting operator...")

	cdbClient, err := k8s.NewCockroachDBClient(namespace)
	if err != nil {
		zap.S().Fatalf("Failed to watch cockroach dbs: %+v", err)
	}

	ccClient, err := k8s.NewCockroachClientClient(namespace)
	if err != nil {
		zap.S().Fatalf("Failed to watch cockroach clients: %+v", err)
	}

	cmClient, err := k8s.NewCockroachMigrationClient(namespace)
	if err != nil {
		zap.S().Fatalf("Failed to watch cockroach migrations: %+v", err)
	}

	rdbClient, err := k8s.NewRedisDBClient(namespace)
	if err != nil {
		zap.S().Fatalf("Failed to watch redis dbs: %+v", err)
	}

	m := manager.New(cdbClient, ccClient, cmClient, rdbClient)
	cancel, err := m.Start()
	if err != nil {
		zap.S().Fatalf("Failed to start the manager: %+v", err)
	}

	zap.S().Info("Running!")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")
	cancel()
}
