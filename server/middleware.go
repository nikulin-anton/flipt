package server

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	errs "go.flipt.io/flipt/errors"
	flipt "go.flipt.io/flipt/rpc/flipt"
	"go.flipt.io/flipt/server/cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// ValidationUnaryInterceptor validates incomming requests
func ValidationUnaryInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if v, ok := req.(flipt.Validator); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

// ErrorUnaryInterceptor intercepts known errors and returns the appropriate GRPC status code
func ErrorUnaryInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	errorsTotal.Inc()

	var errnf errs.ErrNotFound
	if errors.As(err, &errnf) {
		err = status.Error(codes.NotFound, err.Error())
		return
	}

	var errin errs.ErrInvalid
	if errors.As(err, &errin) {
		err = status.Error(codes.InvalidArgument, err.Error())
		return
	}

	var errv errs.ErrValidation
	if errors.As(err, &errv) {
		err = status.Error(codes.InvalidArgument, err.Error())
		return
	}

	err = status.Error(codes.Internal, err.Error())
	return
}

func flagCachingUnaryInterceptor(cache cache.Cacher, logger logrus.FieldLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if cache == nil {
			return handler(ctx, req)
		}

		switch r := req.(type) {
		case *flipt.GetFlagRequest:
			key := flagCacheKey(r.GetKey())

			cached, ok, err := cache.Get(ctx, key)
			if err != nil {
				// if error, log and continue without cache
				logger.WithError(err).Error("getting from cache")
			}

			if ok {
				// if cached, return it
				flag := &flipt.Flag{}
				if err := proto.Unmarshal(cached, flag); err != nil {
					logger.WithError(err).Error("unmarshalling from cache")
					return handler(ctx, req)
				}

				logger.Debugf("flag cache hit: %+v", flag)
				return flag, nil
			}

			logger.Debug("flag cache miss")
			resp, err := handler(ctx, req)
			if err != nil {
				return nil, err
			}

			// marshal response
			data, merr := proto.Marshal(resp.(*flipt.Flag))
			if merr != nil {
				logger.WithError(merr).Error("marshalling for cache")
				return resp, err
			}

			// set in cache
			if cerr := cache.Set(ctx, key, data); cerr != nil {
				logger.WithError(cerr).Error("setting in cache")
			}

			return resp, err

		case *flipt.UpdateFlagRequest, *flipt.DeleteFlagRequest:
			// need to do this assertion because the request type is not known in this block
			keyer := r.(flagKeyer)
			// delete from cache
			if err := cache.Delete(ctx, flagCacheKey(keyer.GetKey())); err != nil {
				logger.WithError(err).Error("deleting from cache")
			}
		case *flipt.CreateVariantRequest, *flipt.UpdateVariantRequest, *flipt.DeleteVariantRequest:
			// need to do this assertion because the request type is not known in this block
			keyer := r.(variantFlagKeyger)
			// delete from cache
			if err := cache.Delete(ctx, flagCacheKey(keyer.GetFlagKey())); err != nil {
				logger.WithError(err).Error("deleting from cache")
			}
		}

		return handler(ctx, req)
	}
}

type flagKeyer interface {
	GetKey() string
}

type variantFlagKeyger interface {
	GetFlagKey() string
}

func flagCacheKey(key string) string {
	k := fmt.Sprintf("f:%s", key)
	return fmt.Sprintf("flipt:%x", md5.Sum([]byte(k))) //nolint:gosec
}
