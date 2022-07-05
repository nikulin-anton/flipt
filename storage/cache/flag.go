package cache

import (
	"context"

	flipt "go.flipt.io/flipt/rpc/flipt"
	"go.flipt.io/flipt/storage"
)

const flagCachePrefix = "flag:"

// GetFlag returns the flag from the cache if it exists; otherwise it delegates to the underlying store
// caching the result if no error
func (c *Store) GetFlag(ctx context.Context, k string) (*flipt.Flag, error) {
	var (
		key   = flagCachePrefix + k
		cache = c.cache.String()
	)

	// check if flag exists in cache
	data, ok, _ := c.cache.Get(ctx, key)

	if ok {
		c.logger.Debugf("cache hit: %q", key)

		flag, ok := data.(*flipt.Flag)
		if !ok {
			// not flag, bad cache
			_ = c.cache.Delete(ctx, key)
			return nil, ErrCorrupt
		}

		return flag, nil
	}

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
