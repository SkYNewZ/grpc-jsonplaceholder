package main

import (
	"context"
	"net/http"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/SkYNewZ/grpc-jsonplaceholder/internal/genproto"
)

func (s userService) GetUser(ctx context.Context, request *genproto.GetUserRequest) (*genproto.User, error) {
	var id = strconv.Itoa(int(request.GetId()))
	var result genproto.User
	if err := s.fetch(ctx, http.MethodGet, s.serverURL+"/users/"+id, nil, &result); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &result, nil
}
