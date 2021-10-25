package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/SkYNewZ/grpc-jsonplaceholder/internal/genproto"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type commentService struct {
	genproto.UnimplementedCommentServiceServer
	healthpb.UnimplementedHealthServer
	client    *http.Client
	serverURL string

	postServiceAddr string
	postServiceConn *grpc.ClientConn
}

func newCommentService() *commentService {
	s := &commentService{
		client:    new(http.Client),
		serverURL: "https://jsonplaceholder.typicode.com",
	}

	// register our clients
	ctx := context.Background()

	mustMapEnv(&s.postServiceAddr, "POST_SERVICE_ADDR")
	mustConnGRPC(ctx, &s.postServiceConn, s.postServiceAddr, grpc.WithInsecure())

	return s
}

func (c commentService) fetch(ctx context.Context, method, url string, body io.Reader, data interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext: %v", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("http.Client.Do: %v", err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		err = fmt.Errorf("json.Decoder.Decode: %v", err)
	}

	return err
}

func (c commentService) GetComment(ctx context.Context, request *genproto.GetCommentRequest) (*genproto.Comment, error) {
	return c.UnimplementedCommentServiceServer.GetComment(ctx, request)
}

func (c commentService) ListComments(ctx context.Context, empty *genproto.Empty) (*genproto.ListCommentsResponse, error) {
	return c.UnimplementedCommentServiceServer.ListComments(ctx, empty)
}
