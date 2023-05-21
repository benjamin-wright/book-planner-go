package events

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func Run(handlers []HandlerFunc) error {
	wg := sync.WaitGroup{}
	wg.Add(1)

	nc, err := nats.Connect(
		os.Getenv("EVENTS_URL"),
		nats.ClosedHandler(func(c *nats.Conn) {
			wg.Done()
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to NATs event bus: %+v", err)
	}

	c, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	if err != nil {
		return fmt.Errorf("failed to create GOB-encoded connection: %+v", err)
	}

	for _, handler := range handlers {
		_, err := handler(c)
		if err != nil {
			return fmt.Errorf("failed to subscribe to book events subject: %+v", err)
		}
	}

	zap.S().Info("Running")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit

	zap.S().Info("Received interrupt or SIGTERM, draining queue")

	nc.Drain()
	wg.Wait()

	zap.S().Info("Finished")
	return nil
}
