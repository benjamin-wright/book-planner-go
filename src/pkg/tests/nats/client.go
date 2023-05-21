package nats

import (
	"os"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/events"
)

func Connect(t *testing.T) *nats.EncodedConn {
	nc, err := nats.Connect(os.Getenv("EVENTS_URL"))
	if err != nil {
		t.Errorf("failed to connect to NATs event bus: %+v", err)
		t.FailNow()
	}

	c, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	if err != nil {
		t.Errorf("failed to create GOB-encoded connection: %+v", err)
		t.FailNow()
	}

	return c
}

func Handle(t *testing.T, nc *nats.EncodedConn, handler events.HandlerFunc) {
	_, err := handler(nc)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
}
