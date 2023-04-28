package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/operators/db/pkg/k8s"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	namespace := os.Getenv("NAMESPACE")

	zap.S().Info("Starting operator...")

	cli, err := k8s.New()
	if err != nil {
		zap.S().Fatalf("Failed to create k8s client: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	dbs, err := cli.CockroachDBWatch(ctx, cancel, namespace)
	if err != nil {
		zap.S().Fatalf("Failed to watch cockroachdbs: %+v", err)
	}

	go func(dbs <-chan map[string]*k8s.CockroachDB) {
		for db := range dbs {
			zap.S().Infof("Event: %+v", db)
		}
	}(dbs)

	zap.S().Info("Running!")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")
	cancel()
}
