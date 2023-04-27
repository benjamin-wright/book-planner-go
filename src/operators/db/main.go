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

	zap.S().Info("Starting operator...")

	cli, err := k8s.New()
	if err != nil {
		zap.S().Fatalf("Failed to create k8s client: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cli.WatchCockroachDBs(ctx, cancel, func(old k8s.CockroachDB, new k8s.CockroachDB) {})

	zap.S().Info("Running!")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")
	cancel()
}
