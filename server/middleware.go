package server

import (
	"context"
	"errors"

	errs "go.flipt.io/flipt/errors"
	flipt "go.flipt.io/flipt/rpc/flipt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
