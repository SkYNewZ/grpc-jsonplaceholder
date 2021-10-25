package main

import (
	"context"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func (c commentService) Check(ctx context.Context, request *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return c.UnimplementedHealthServer.Check(ctx, request)
}

func (c commentService) Watch(request *healthpb.HealthCheckRequest, server healthpb.Health_WatchServer) error {
	return c.UnimplementedHealthServer.Watch(request, server)
}
