package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.flipt.io/flipt/config"
	flipt "go.flipt.io/flipt/rpc/flipt"
	"go.flipt.io/flipt/server/cache/memory"
)

func TestGetFlag(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
		req = &flipt.GetFlagRequest{Key: "foo"}
	)

	store.On("GetFlag", mock.Anything, "foo").Return(&flipt.Flag{
		Key:     req.Key,
		Enabled: true,
	}, nil)

	got, err := s.GetFlag(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)
	assert.Equal(t, "foo", got.Key)
	assert.Equal(t, true, got.Enabled)
}

func TestGetFlag_WithCache(t *testing.T) {
	var (
		store = &storeMock{}
		cache = memory.NewCache(config.CacheConfig{
			TTL:     time.Second,
			Enabled: true,
			Backend: config.CacheMemory,
		})
		cacheSpy = newCacheSpy(cache)
		s        = &Server{
			logger:       logger,
			store:        store,
			cache:        cacheSpy,
			cacheEnabled: true,
		}
		req = &flipt.GetFlagRequest{Key: "foo"}
	)

	store.On("GetFlag", mock.Anything, "foo").Return(&flipt.Flag{
		Key:     req.Key,
		Enabled: true,
	}, nil)

	// run many times to ensure cache is working correctly
	for i := 0; i < 10; i++ {
		got, err := s.GetFlag(context.TODO(), req)
		require.NoError(t, err)

		assert.NotNil(t, got)
		assert.Equal(t, "foo", got.Key)
		assert.Equal(t, true, got.Enabled)
	}

	assert.Equal(t, 10, cacheSpy.getCalled)
	assert.NotEmpty(t, cacheSpy.getKeys)

	// cache key is flipt:(md5(f:foo))
	const cacheKey = "flipt:864ce319cc64891a59e4745fbe7ecc47"
	_, ok := cacheSpy.getKeys[cacheKey]
	assert.True(t, ok)

	assert.Equal(t, 1, cacheSpy.setCalled)
	assert.NotEmpty(t, cacheSpy.setItems)
	assert.NotEmpty(t, cacheSpy.setItems[cacheKey])

}

func TestListFlags(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
	)

	store.On("ListFlags", mock.Anything, mock.Anything).Return(
		[]*flipt.Flag{
			{
				Key: "foo",
			},
		}, nil)

	got, err := s.ListFlags(context.TODO(), &flipt.ListFlagRequest{})
	require.NoError(t, err)

	assert.NotEmpty(t, got.Flags)
}

func TestListFlags_NoCache(t *testing.T) {
	var (
		store = &storeMock{}
		cache = memory.NewCache(config.CacheConfig{
			TTL:     time.Second,
			Enabled: true,
			Backend: config.CacheMemory,
		})
		cacheSpy = newCacheSpy(cache)
		s        = &Server{
			logger:       logger,
			store:        store,
			cache:        cacheSpy,
			cacheEnabled: true,
		}
	)

	store.On("ListFlags", mock.Anything, mock.Anything).Return(
		[]*flipt.Flag{
			{
				Key: "foo",
			},
		}, nil)

	got, err := s.ListFlags(context.TODO(), &flipt.ListFlagRequest{})
	require.NoError(t, err)

	assert.NotEmpty(t, got.Flags)

	assert.Equal(t, 0, cacheSpy.getCalled)
	assert.Empty(t, cacheSpy.getKeys)

	assert.Equal(t, 0, cacheSpy.setCalled)
	assert.Empty(t, cacheSpy.setItems)
}

func TestCreateFlag(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
		req = &flipt.CreateFlagRequest{
			Key:         "key",
			Name:        "name",
			Description: "desc",
			Enabled:     true,
		}
	)

	store.On("CreateFlag", mock.Anything, req).Return(&flipt.Flag{
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
	}, nil)

	got, err := s.CreateFlag(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)
}

func TestCreateFlag_NoCache(t *testing.T) {
	var (
		store = &storeMock{}
		cache = memory.NewCache(config.CacheConfig{
			TTL:     time.Second,
			Enabled: true,
			Backend: config.CacheMemory,
		})
		cacheSpy = newCacheSpy(cache)
		s        = &Server{
			logger:       logger,
			store:        store,
			cache:        cacheSpy,
			cacheEnabled: true,
		}
		req = &flipt.CreateFlagRequest{
			Key:         "key",
			Name:        "name",
			Description: "desc",
			Enabled:     true,
		}
	)

	store.On("CreateFlag", mock.Anything, req).Return(&flipt.Flag{
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
	}, nil)

	got, err := s.CreateFlag(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)

	assert.Equal(t, 0, cacheSpy.getCalled)
	assert.Empty(t, cacheSpy.getKeys)

	assert.Equal(t, 0, cacheSpy.setCalled)
	assert.Empty(t, cacheSpy.setItems)
}

func TestUpdateFlag(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
		req = &flipt.UpdateFlagRequest{
			Key:         "key",
			Name:        "name",
			Description: "desc",
			Enabled:     true,
		}
	)

	store.On("UpdateFlag", mock.Anything, req).Return(&flipt.Flag{
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
	}, nil)

	got, err := s.UpdateFlag(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)
}

func TestUpdateFlag_WithCache(t *testing.T) {
	var (
		store = &storeMock{}
		cache = memory.NewCache(config.CacheConfig{
			TTL:     time.Second,
			Enabled: true,
			Backend: config.CacheMemory,
		})
		cacheSpy = newCacheSpy(cache)
		s        = &Server{
			logger:       logger,
			store:        store,
			cache:        cacheSpy,
			cacheEnabled: true,
		}
		req = &flipt.UpdateFlagRequest{
			Key:         "key",
			Name:        "name",
			Description: "desc",
			Enabled:     true,
		}
	)

	store.On("UpdateFlag", mock.Anything, req).Return(&flipt.Flag{
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
	}, nil)

	got, err := s.UpdateFlag(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)

	assert.Equal(t, 1, cacheSpy.deleteCalled)
	assert.NotEmpty(t, cacheSpy.deleteKeys)
}

func TestDeleteFlag(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
		req = &flipt.DeleteFlagRequest{
			Key: "key",
		}
	)

	store.On("DeleteFlag", mock.Anything, req).Return(nil)

	got, err := s.DeleteFlag(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)
}

func TestDeleteFlag_WithCache(t *testing.T) {
	var (
		store = &storeMock{}
		cache = memory.NewCache(config.CacheConfig{
			TTL:     time.Second,
			Enabled: true,
			Backend: config.CacheMemory,
		})
		cacheSpy = newCacheSpy(cache)
		s        = &Server{
			logger:       logger,
			store:        store,
			cache:        cacheSpy,
			cacheEnabled: true,
		}
		req = &flipt.DeleteFlagRequest{
			Key: "key",
		}
	)

	store.On("DeleteFlag", mock.Anything, req).Return(nil)

	got, err := s.DeleteFlag(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)

	assert.Equal(t, 1, cacheSpy.deleteCalled)
	assert.NotEmpty(t, cacheSpy.deleteKeys)
}

func TestCreateVariant(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
		req = &flipt.CreateVariantRequest{
			FlagKey:     "flagKey",
			Key:         "key",
			Name:        "name",
			Description: "desc",
		}
	)

	store.On("CreateVariant", mock.Anything, req).Return(&flipt.Variant{
		Id:          "1",
		FlagKey:     req.FlagKey,
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Attachment:  req.Attachment,
	}, nil)

	got, err := s.CreateVariant(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)
}

func TestCreateVariant_WithCache(t *testing.T) {
	var (
		store = &storeMock{}
		cache = memory.NewCache(config.CacheConfig{
			TTL:     time.Second,
			Enabled: true,
			Backend: config.CacheMemory,
		})
		cacheSpy = newCacheSpy(cache)
		s        = &Server{
			logger:       logger,
			store:        store,
			cache:        cacheSpy,
			cacheEnabled: true,
		}
		req = &flipt.CreateVariantRequest{
			FlagKey:     "flagKey",
			Key:         "key",
			Name:        "name",
			Description: "desc",
		}
	)

	store.On("CreateVariant", mock.Anything, req).Return(&flipt.Variant{
		Id:          "1",
		FlagKey:     req.FlagKey,
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Attachment:  req.Attachment,
	}, nil)

	got, err := s.CreateVariant(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)

	assert.Equal(t, 1, cacheSpy.deleteCalled)
	assert.NotEmpty(t, cacheSpy.deleteKeys)
}

func TestUpdateVariant(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
		req = &flipt.UpdateVariantRequest{
			Id:          "1",
			FlagKey:     "flagKey",
			Key:         "key",
			Name:        "name",
			Description: "desc",
		}
	)

	store.On("UpdateVariant", mock.Anything, req).Return(&flipt.Variant{
		Id:          req.Id,
		FlagKey:     req.FlagKey,
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Attachment:  req.Attachment,
	}, nil)

	got, err := s.UpdateVariant(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)
}

func TestUpdateVariant_WithCache(t *testing.T) {
	var (
		store = &storeMock{}
		cache = memory.NewCache(config.CacheConfig{
			TTL:     time.Second,
			Enabled: true,
			Backend: config.CacheMemory,
		})
		cacheSpy = newCacheSpy(cache)
		s        = &Server{
			logger:       logger,
			store:        store,
			cache:        cacheSpy,
			cacheEnabled: true,
		}
		req = &flipt.UpdateVariantRequest{
			Id:          "1",
			FlagKey:     "flagKey",
			Key:         "key",
			Name:        "name",
			Description: "desc",
		}
	)

	store.On("UpdateVariant", mock.Anything, req).Return(&flipt.Variant{
		Id:          req.Id,
		FlagKey:     req.FlagKey,
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Attachment:  req.Attachment,
	}, nil)

	got, err := s.UpdateVariant(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)

	assert.Equal(t, 1, cacheSpy.deleteCalled)
	assert.NotEmpty(t, cacheSpy.deleteKeys)
}

func TestDeleteVariant(t *testing.T) {
	var (
		store = &storeMock{}
		s     = &Server{
			logger: logger,
			store:  store,
		}
		req = &flipt.DeleteVariantRequest{
			Id: "1",
		}
	)

	store.On("DeleteVariant", mock.Anything, req).Return(nil)

	got, err := s.DeleteVariant(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)
}

func TestDeleteVariant_WithCache(t *testing.T) {
	var (
		store = &storeMock{}
		cache = memory.NewCache(config.CacheConfig{
			TTL:     time.Second,
			Enabled: true,
			Backend: config.CacheMemory,
		})
		cacheSpy = newCacheSpy(cache)
		s        = &Server{
			logger:       logger,
			store:        store,
			cache:        cacheSpy,
			cacheEnabled: true,
		}
		req = &flipt.DeleteVariantRequest{
			Id: "1",
		}
	)

	store.On("DeleteVariant", mock.Anything, req).Return(nil)

	got, err := s.DeleteVariant(context.TODO(), req)
	require.NoError(t, err)

	assert.NotNil(t, got)

	assert.Equal(t, 1, cacheSpy.deleteCalled)
	assert.NotEmpty(t, cacheSpy.deleteKeys)
}
