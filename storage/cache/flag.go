package cache

import (
	"context"

	flipt "github.com/markphelps/flipt/rpc/flipt"
)

const flagCachePrefix = "flag:"

// GetFlag returns the flag from the cache if it exists; otherwise it delegates to the underlying store
// caching the result if no error
func (c *Store) GetFlag(ctx context.Context, k string) (*flipt.Flag, error) {
	key := flagCachePrefix + k

	// check if flag exists in cache
	data, ok, _ := c.cache.Get(ctx, key)

	if ok {
		c.logger.Debugf("cache hit: %q", key)

		flag, ok := data.(*flipt.Flag)
		if !ok {
			// not flag, bad cache
			c.logger.Errorf("corrupt cache, deleting: %q", key)
			if err := c.cache.Delete(ctx, key); err != nil {
				c.logger.WithError(err).Error("deleting cache entry")
			}
			goto db
		}

		return flag, nil
	}

db:
	// flag not in cache, delegate to underlying store
	flag, err := c.Store.GetFlag(ctx, k)
	if err != nil {
		return flag, err
	}

	_ = c.cache.Set(ctx, key, flag)

	c.logger.Debugf("cache miss; added: %q", key)
	return flag, nil
}

// UpdateFlag delegates to the underlying store, deleting flag from the cache in the process
func (c *Store) UpdateFlag(ctx context.Context, r *flipt.UpdateFlagRequest) (*flipt.Flag, error) {
	key := flagCachePrefix + r.Key
	_ = c.cache.Delete(ctx, key)
	c.logger.Debugf("deleted flag from cache: %q", key)
	return c.Store.UpdateFlag(ctx, r)
}

// DeleteFlag delegates to the underlying store, deleting flag from the cache in the process
func (c *Store) DeleteFlag(ctx context.Context, r *flipt.DeleteFlagRequest) error {
	key := flagCachePrefix + r.Key
	_ = c.cache.Delete(ctx, key)
	c.logger.Debugf("deleted flag from cache: %q", key)
	return c.Store.DeleteFlag(ctx, r)
}

// CreateVariant delegates to the underlying store, deleting flag from the cache in the process
func (c *Store) CreateVariant(ctx context.Context, r *flipt.CreateVariantRequest) (*flipt.Variant, error) {
	key := flagCachePrefix + r.FlagKey
	_ = c.cache.Delete(ctx, key)
	c.logger.Debugf("deleted flag from cache: %q", key)
	return c.Store.CreateVariant(ctx, r)
}

// UpdateVariant delegates to the underlying store, deleting flag from the cache in the process
func (c *Store) UpdateVariant(ctx context.Context, r *flipt.UpdateVariantRequest) (*flipt.Variant, error) {
	key := flagCachePrefix + r.FlagKey
	_ = c.cache.Delete(ctx, key)
	c.logger.Debugf("deleted flag from cache: %q", key)
	return c.Store.UpdateVariant(ctx, r)
}

// DeleteVariant delegates to the underlying store, deleting flag from the cache in the process
func (c *Store) DeleteVariant(ctx context.Context, r *flipt.DeleteVariantRequest) error {
	key := flagCachePrefix + r.FlagKey
	_ = c.cache.Delete(ctx, key)
	c.logger.Debugf("deleted flag from cache: %q", key)
	return c.Store.DeleteVariant(ctx, r)
}
