package cache

import (
	"context"
	"fmt"
)

// Cacher modifies and queries a cache
type Cacher interface {
	// Get retrieves a value from the cache, the bool indicates if the item was found
	Get(ctx context.Context, key string) (interface{}, bool, error)
	// Set sets a value in the cache
	Set(ctx context.Context, key string, value interface{}) error
	// Delete removes a value from the cache
	Delete(ctx context.Context, key string) error
	fmt.Stringer
}
