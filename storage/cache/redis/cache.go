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

func (r *Redis) Get(key string) (i interface{}, b bool) {
	// TODO: add context to cacher interface
	ctx := context.Background()
	if err := r.c.Get(ctx, key, i); err != nil {
		return nil, false
	}
	return i, true
}

func (r *Redis) Set(key string, value interface{}) {
	// TODO: add context to cacher interface
	ctx := context.Background()
	// TODO: modify cacher.Set to return err
	_ = r.c.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
	})
}

func (r *Redis) Delete(key string) {
	panic("not implemented") // TODO: Implement
}

func (r *Redis) Flush() {
	panic("not implemented") // TODO: Implement
}

func (r *Redis) String() string {
	return "redis"
}
