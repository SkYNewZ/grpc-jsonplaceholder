package main

import (
	"context"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func (p postService) Check(ctx context.Context, request *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return p.UnimplementedHealthServer.Check(ctx, request)
}

func (p postService) Watch(request *healthpb.HealthCheckRequest, server healthpb.Health_WatchServer) error {
	return p.UnimplementedHealthServer.Watch(request, server)
}
