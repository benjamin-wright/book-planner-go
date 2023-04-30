package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/operators/db/pkg/k8s"
	"ponglehub.co.uk/book-planner-go/src/operators/db/pkg/manager"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	namespace := os.Getenv("NAMESPACE")

	zap.S().Info("Starting operator...")

	ctx, cancel := context.WithCancel(context.Background())

	cdbs, err := k8s.WatchCockroachDBs(ctx, cancel, namespace)
	if err != nil {
		zap.S().Fatalf("Failed to watch cockroach dbs: %+v", err)
	}

	cclients, err := k8s.WatchCockroachClients(ctx, cancel, namespace)
	if err != nil {
		zap.S().Fatalf("Failed to watch cockroach clients: %+v", err)
	}

	cmigrations, err := k8s.WatchCockroachMigrations(ctx, cancel, namespace)
	if err != nil {
		zap.S().Fatalf("Failed to watch cockroach migrations: %+v", err)
	}

	rdbs, err := k8s.WatchRedisDBs(ctx, cancel, namespace)
	if err != nil {
		zap.S().Fatalf("Failed to watch redis dbs: %+v", err)
	}

	manager.Manage(ctx, cdbs, cclients, cmigrations, rdbs)

	zap.S().Info("Running!")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")
	cancel()
}
