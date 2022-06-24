package server

import (
	"context"
	"crypto/md5"
	"fmt"

	flipt "go.flipt.io/flipt/rpc/flipt"
	"go.flipt.io/flipt/storage"
	"google.golang.org/protobuf/proto"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

// GetFlag gets a flag
func (s *Server) GetFlag(ctx context.Context, r *flipt.GetFlagRequest) (flag *flipt.Flag, err error) {
	s.logger.WithField("request", r).Debug("get flag")
	if s.cacheEnabled {
		// check cache
		flag, err = s.getFlagWithCache(ctx, r)
	} else {
		flag, err = s.store.GetFlag(ctx, r.Key)
	}
	s.logger.WithField("response", flag).Debug("get flag")
	return flag, err
}

func (s *Server) getFlagWithCache(ctx context.Context, r *flipt.GetFlagRequest) (*flipt.Flag, error) {
	var (
		logger   = s.logger.WithField("request", r)
		key, err = flagCacheKey(r)
	)

	if err != nil {
		// if error, log and continue without cache
		logger.WithError(err).Error("generating cache key")
		return s.store.GetFlag(ctx, r.Key)
	}

	cached, ok, err := s.cache.Get(ctx, key)
	if err != nil {
		// if error, log and continue without cache
		logger.WithError(err).Error("getting from cache")
		return s.store.GetFlag(ctx, r.Key)
	}

	if !ok {
		logger.Debug("flag cache miss")
		flag, err := s.store.GetFlag(ctx, r.Key)
		if err != nil {
			return flag, err
		}
		data, err := proto.Marshal(flag)
		if err != nil {
			return flag, err
		}
		err = s.cache.Set(ctx, key, data)
		return flag, err
	}

	flag := &flipt.Flag{}
	if err := proto.Unmarshal(cached, flag); err != nil {
		logger.WithError(err).Error("unmarshalling from cache")
		return s.store.GetFlag(ctx, r.Key)
	}

	logger.Debugf("flag cache hit: %+v", flag)
	return flag, nil
}

// ListFlags lists all flags
func (s *Server) ListFlags(ctx context.Context, r *flipt.ListFlagRequest) (*flipt.FlagList, error) {
	s.logger.WithField("request", r).Debug("list flags")

	flags, err := s.store.ListFlags(ctx, storage.WithLimit(uint64(r.Limit)), storage.WithOffset(uint64(r.Offset)))
	if err != nil {
		return nil, err
	}

	var resp flipt.FlagList

	for i := range flags {
		resp.Flags = append(resp.Flags, flags[i])
	}

	s.logger.WithField("response", &resp).Debug("list flags")
	return &resp, nil
}

// CreateFlag creates a flag
func (s *Server) CreateFlag(ctx context.Context, r *flipt.CreateFlagRequest) (*flipt.Flag, error) {
	s.logger.WithField("request", r).Debug("create flag")
	flag, err := s.store.CreateFlag(ctx, r)
	s.logger.WithField("response", flag).Debug("create flag")
	return flag, err
}

// UpdateFlag updates an existing flag
func (s *Server) UpdateFlag(ctx context.Context, r *flipt.UpdateFlagRequest) (*flipt.Flag, error) {
	s.logger.WithField("request", r).Debug("update flag")
	flag, err := s.store.UpdateFlag(ctx, r)
	s.logger.WithField("response", flag).Debug("update flag")
	return flag, err
}

// DeleteFlag deletes a flag
func (s *Server) DeleteFlag(ctx context.Context, r *flipt.DeleteFlagRequest) (*empty.Empty, error) {
	s.logger.WithField("request", r).Debug("delete flag")
	if err := s.store.DeleteFlag(ctx, r); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// CreateVariant creates a variant
func (s *Server) CreateVariant(ctx context.Context, r *flipt.CreateVariantRequest) (*flipt.Variant, error) {
	s.logger.WithField("request", r).Debug("create variant")
	variant, err := s.store.CreateVariant(ctx, r)
	s.logger.WithField("response", variant).Debug("create variant")
	return variant, err
}

// UpdateVariant updates an existing variant
func (s *Server) UpdateVariant(ctx context.Context, r *flipt.UpdateVariantRequest) (*flipt.Variant, error) {
	s.logger.WithField("request", r).Debug("update variant")
	variant, err := s.store.UpdateVariant(ctx, r)
	s.logger.WithField("response", variant).Debug("update variant")
	return variant, err
}

// DeleteVariant deletes a variant
func (s *Server) DeleteVariant(ctx context.Context, r *flipt.DeleteVariantRequest) (*empty.Empty, error) {
	s.logger.WithField("request", r).Debug("delete variant")
	if err := s.store.DeleteVariant(ctx, r); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func flagCacheKey(r *flipt.GetFlagRequest) (string, error) {
	k := fmt.Sprintf("f:%s", r.GetKey())
	return fmt.Sprintf("flipt:%x", md5.Sum([]byte(k))), nil //nolint:gosec
}
