package main

import (
	"context"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func (s userService) Check(ctx context.Context, request *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return s.UnimplementedHealthServer.Check(ctx, request)
}

func (s userService) Watch(request *healthpb.HealthCheckRequest, server healthpb.Health_WatchServer) error {
	return s.UnimplementedHealthServer.Watch(request, server)
}
