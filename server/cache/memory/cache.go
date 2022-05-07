package memory

import (
	"context"
	"sync"
	"time"

	"github.com/markphelps/flipt/server/cache"
	gocache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

// InMemoryCache wraps gocache.Cache in order to implement Cacher
type InMemoryCache struct {
	c             *gocache.Cache
	mu            sync.RWMutex
	itemCount     int64
	missTotal     int64
	hitTotal      int64
	evictionTotal int64
}

// NewCache creates a new InMemoryCache with the provided ttl and evictionInterval
func NewCache(ttl time.Duration, evictionInterval time.Duration, logger logrus.FieldLogger) *InMemoryCache {
	logger = logger.WithField("cache", "memory")

	var (
		c     = gocache.New(ttl, evictionInterval)
		cache = &InMemoryCache{c: c}
	)

	c.OnEvicted(func(s string, _ interface{}) {
		cache.mu.Lock()
		cache.itemCount--
		cache.evictionTotal++
		cache.mu.Unlock()
		logger.Debugf("evicted key: %q", s)
	})

	return cache
}

func (i *InMemoryCache) Get(_ context.Context, key string) (interface{}, bool, error) {
	v, ok := i.c.Get(key)
	if !ok {
		i.mu.Lock()
		i.missTotal++
		i.mu.Unlock()
		return nil, false, nil
	}

	i.mu.Lock()
	i.hitTotal++
	i.mu.Unlock()
	return v, true, nil
}

func (i *InMemoryCache) Set(_ context.Context, key string, value interface{}) error {
	i.c.SetDefault(key, value)
	i.mu.Lock()
	i.itemCount++
	i.mu.Unlock()
	return nil
}

func (i *InMemoryCache) Delete(_ context.Context, key string) error {
	i.c.Delete(key)
	i.mu.Lock()
	i.itemCount--
	i.mu.Unlock()
	return nil
}

func (i *InMemoryCache) String() string {
	return "memory"
}

func (i *InMemoryCache) Stats() cache.Stats {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return cache.Stats{
		ItemCount:     i.itemCount,
		MissTotal:     i.missTotal,
		HitTotal:      i.hitTotal,
		EvictionTotal: i.evictionTotal,
	}
}
