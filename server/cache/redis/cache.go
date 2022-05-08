package redis

import (
	"context"
	"errors"
	"sync/atomic"

	redis "github.com/go-redis/cache/v8"
	"github.com/markphelps/flipt/config"
	"github.com/markphelps/flipt/server/cache"
)

type Cache struct {
	c         *redis.Cache
	cfg       config.CacheConfig
	missTotal uint64
	hitTotal  uint64
	errTotal  uint64
}

// NewCache creates a new redis cache with the provided cache config
func NewCache(cfg config.CacheConfig, c *redis.Cache) *Cache {
	return &Cache{cfg: cfg, c: c}
}

func (c *Cache) Get(ctx context.Context, key string) (i interface{}, b bool, err error) {
	if err := c.c.Get(ctx, key, i); err != nil {
		if errors.Is(err, redis.ErrCacheMiss) {
			atomic.AddUint64(&c.missTotal, 1)
			return nil, false, nil
		}

		atomic.AddUint64(&c.errTotal, 1)
		return nil, false, err
	}

	atomic.AddUint64(&c.hitTotal, 1)
	return i, true, nil
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}) error {
	if err := c.c.Set(&redis.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   c.cfg.TTL,
	}); err != nil {
		atomic.AddUint64(&c.errTotal, 1)
		return err
	}

	return nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	if err := c.c.Delete(ctx, key); err != nil {
		atomic.AddUint64(&c.errTotal, 1)
		return err
	}

	return nil
}

func (c *Cache) String() string {
	return "redis"
}

func (c *Cache) Stats() cache.Stats {
	return cache.Stats{
		MissTotal:  c.missTotal,
		HitTotal:   c.hitTotal,
		ErrorTotal: c.errTotal,
	}
}
