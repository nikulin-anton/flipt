package server

import (
	"context"
	"errors"

	errs "go.flipt.io/flipt/errors"
	flipt "go.flipt.io/flipt/rpc/flipt"
	"go.flipt.io/flipt/storage"

	"github.com/sirupsen/logrus"
)

var _ flipt.FliptServer = &Server{}

type Option func(s *Server)

// Server serves the Flipt backend
type Server struct {
	logger logrus.FieldLogger
	store  storage.Store

	cacheEnabled bool
	cache        cache.Cacher

	flipt.UnimplementedFliptServer
}

// New creates a new Server
func New(logger logrus.FieldLogger, store storage.Store, opts ...Option) *Server {
	var (
		s = &Server{
			logger: logger,
			store:  store,
		}
	)

	for _, fn := range opts {
		fn(s)
	}

	return s
}

func WithCache(cache cache.Cacher) Option {
	return func(s *Server) {
		s.cache = cache
		s.cacheEnabled = true
	}
}
