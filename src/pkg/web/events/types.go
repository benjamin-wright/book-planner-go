package events

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type HandlerFunc func(c *nats.EncodedConn) (*nats.Subscription, error)

type RequestHandler[T any, U any] struct {
	Channel string
	Handler func(event *T) U
}

func (h RequestHandler[T, U]) GetHandler() HandlerFunc {
	return func(c *nats.EncodedConn) (*nats.Subscription, error) {
		zap.S().Infof("Subscribing to '%s' events", h.Channel)

		return c.Subscribe(h.Channel, func(subject, reply string, event *T) {
			if event == nil {
				zap.S().Error("Received nil event")
				return
			}

			response := h.Handler(event)

			err := c.Publish(reply, &response)
			if err != nil {
				zap.S().Errorf("Failed to send response: %+v", err)
			}
		})
	}
}
