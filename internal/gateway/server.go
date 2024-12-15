package gateway

import (
	api "github.com/wanderer69/user_registration/pkg/api/public"
)

type Server struct {
	authOperations authOperations
}

var _ api.ServerInterface = (*Server)(nil)

func NewServer(
	authOperations authOperations,
) *Server {
	return &Server{
		authOperations: authOperations,
	}
}
