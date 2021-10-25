package main

import (
	"context"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/SkYNewZ/grpc-jsonplaceholder/internal/genproto"
)

func (s userService) ListUsers(ctx context.Context, _ *genproto.Empty) (*genproto.ListUsersResponse, error) {
	var results genproto.ListUsersResponse
	if err := s.fetch(ctx, http.MethodGet, s.serverURL+"/users", nil, &results.Users); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &results, nil
}
