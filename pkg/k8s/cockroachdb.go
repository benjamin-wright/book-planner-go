package k8s

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type CockroachDB struct {
	Name string
}

type CockroachDBWatchHandler func(old CockroachDB, new CockroachDB)

func (c *Client) WatchCockroachDBs(ctx context.Context, cancel context.CancelFunc, handler CockroachDBWatchHandler) error {
	watcher, err := c.client.Resource(CockroachDBSchema).Watch(ctx, v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to watch cockroach dbs: %+v", err)
	}

	go func(w <-chan watch.Event) {
		for e := range w {
			switch e.Type {
			case watch.Added:
				{
					zap.S().Debug("Watch event: CockroachDB[Added]")
				}
			case watch.Modified:
				{
					zap.S().Debug("Watch event: CockroachDB[Modified]")
				}
			case watch.Bookmark:
				{
					zap.S().Debug("Watch event: CockroachDB[Bookmark]")
				}
			case watch.Deleted:
				{
					zap.S().Debug("Watch event: CockroachDB[Deleted]")
				}
			case watch.Error:
				{
					zap.S().Debug("Watch event: CockroachDB[Error]")
				}
			}
		}

		cancel()
	}(watcher.ResultChan())

	return nil
}
