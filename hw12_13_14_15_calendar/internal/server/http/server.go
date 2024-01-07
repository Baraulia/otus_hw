package internalhttp

import (
	"context"
)

type Server struct { // TODO
}

type Logger interface { // TODO
}

type Application interface { // TODO
}

//nolint:revive
func NewServer(logger Logger, app Application) *Server {
	return &Server{}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	<-ctx.Done()
	return nil
}

//nolint:revive
func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

// TODO
