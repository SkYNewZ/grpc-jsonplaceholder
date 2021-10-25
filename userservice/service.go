package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/SkYNewZ/grpc-jsonplaceholder/internal/genproto"
)

type userService struct {
	genproto.UnimplementedUserServiceServer
	healthpb.UnimplementedHealthServer
	client    *http.Client
	serverURL string
}

func newUserService() *userService {
	return &userService{
		client:    new(http.Client),
		serverURL: "https://jsonplaceholder.typicode.com",
	}
}

func (s userService) fetch(ctx context.Context, method, url string, body io.Reader, data interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext: %v", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("http.Client.Do: %v", err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		err = fmt.Errorf("json.Decoder.Decode: %v", err)
	}

	return err
}
