package server

import (
	flipt "github.com/markphelps/flipt/rpc/flipt"
	"github.com/markphelps/flipt/server/cache"
	"github.com/markphelps/flipt/storage"

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
