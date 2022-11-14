package infura

import (
	"context"
	"time"

	"interview-test/internal/util/cache"
)

type CachingClient struct {
	cache cache.RefreshingCache[Wei]
}

func NewCachingClient(inner Client, refreshEvery time.Duration, evictionDuration time.Duration) Client {
	return &CachingClient{
		cache: cache.NewRefreshingCache(
			func() (Wei, error) { return inner.GetGasPrice(context.Background()) },
			refreshEvery,
			evictionDuration,
		),
	}
}

func (c CachingClient) GetGasPrice(_ context.Context) (Wei, error) {
	return c.cache.Get()
}
