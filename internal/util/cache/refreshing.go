package cache

import (
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

type RefreshingCache[T any] interface {
	Get() (T, error)
}

type retrievalInfo[T any] struct {
	value         T
	err           error
	lastRetrieval time.Time
}

type refreshingCache[T any] struct {
	cached    atomic.Value
	refreshFn func() (T, error)
}

func (r refreshingCache[T]) Get() (T, error) {
	if info := r.cached.Load().(retrievalInfo[T]); !info.lastRetrieval.IsZero() {
		return info.value, nil
	} else {
		var zero T
		return zero, info.err
	}
}

func (r *refreshingCache[T]) initAndRefreshValue(refreshEvery time.Duration, evictionDuration time.Duration) {
	for {
		w, err := r.refreshFn()

		if err != nil {
			info := r.cached.Load().(retrievalInfo[T])
			info.err = err

			if info.lastRetrieval.Before(time.Now().Add(-evictionDuration)) {
				info.lastRetrieval = time.Time{}
			}

			r.cached.Store(info)
		} else {
			r.cached.Store(retrievalInfo[T]{value: w, lastRetrieval: time.Now()})
		}

		time.Sleep(refreshEvery)
	}
}

func NewRefreshingCache[T any](refreshFn func() (T, error), refreshEvery time.Duration, evictionDuration time.Duration) RefreshingCache[T] {
	c := &refreshingCache[T]{refreshFn: refreshFn}
	c.cached.Store(retrievalInfo[T]{err: errors.New("no cache available")})
	go c.initAndRefreshValue(refreshEvery, evictionDuration)
	return c
}
