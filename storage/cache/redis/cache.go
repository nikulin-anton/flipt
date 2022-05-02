package redis

import (
	"context"

	"github.com/go-redis/cache/v8"
)

type Redis struct {
	c *cache.Cache
}

func NewCache() *Redis {
	return &Redis{}
}

func (r *Redis) Get(ctx context.Context, key string) (i interface{}, b bool, err error) {
	if err := r.c.Get(ctx, key, i); err != nil {
		return nil, false, err
	}
	return i, true, nil
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}) error {
	return r.c.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
	})
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.c.Delete(ctx, key)
}

func (r *Redis) String() string {
	return "redis"
}
