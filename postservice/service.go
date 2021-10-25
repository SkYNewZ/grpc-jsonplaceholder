package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/SkYNewZ/grpc-jsonplaceholder/internal/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type postService struct {
	genproto.UnimplementedPostServiceServer
	healthpb.UnimplementedHealthServer
	client    *http.Client
	serverURL string

	userServiceAddr string
	userServiceConn *grpc.ClientConn
}

func newPostService() *postService {
	s := &postService{
		client:    new(http.Client),
		serverURL: "https://jsonplaceholder.typicode.com",
	}

	// register our clients
	ctx := context.Background()

	mustMapEnv(&s.userServiceAddr, "USER_SERVICE_ADDR")
	mustConnGRPC(ctx, &s.userServiceConn, s.userServiceAddr, grpc.WithInsecure())

	return s
}

func (p postService) fetch(ctx context.Context, method, url string, body io.Reader, data interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext: %v", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("http.Client.Do: %v", err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		err = fmt.Errorf("json.Decoder.Decode: %v", err)
	}

	return err
}

func (p postService) GetPost(ctx context.Context, request *genproto.GetPostRequest) (*genproto.Post, error) {
	type response struct {
		*genproto.Post
		UserID uint32 `json:"userId"`
	}

	var (
		id     = strconv.Itoa(int(request.GetId()))
		err    error
		result response
	)

	if err := p.fetch(ctx, http.MethodGet, p.serverURL+"/posts/"+id, nil, &result); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// fetch user
	client := genproto.NewUserServiceClient(p.userServiceConn)
	result.Post.User, err = client.GetUser(ctx, &genproto.GetUserRequest{Id: result.UserID})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return result.Post, nil
}

func (p postService) ListPosts(ctx context.Context, _ *genproto.Empty) (*genproto.ListPostsResponse, error) {
	type postResponse struct {
		*genproto.Post
		UserID uint32 `json:"userId"`
	}

	type response []*postResponse
	var results response
	if err := p.fetch(ctx, http.MethodGet, p.serverURL+"/posts", nil, &results); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// fetch users
	var (
		wg     sync.WaitGroup
		client = genproto.NewUserServiceClient(p.userServiceConn)
		r      = &genproto.ListPostsResponse{
			Posts: make([]*genproto.Post, len(results)),
		}
	)

	for i, p := range results {
		wg.Add(1)
		go func(post *postResponse, j int) {
			defer wg.Done()

			var err error
			post.Post.User, err = client.GetUser(ctx, &genproto.GetUserRequest{Id: post.UserID})
			if err != nil {
				log.Errorf("client.GetUser: %v", err)
				return
			}

			r.Posts[j] = post.Post
		}(p, i)
	}

	wg.Wait()
	return r, nil
}

func (p postService) CreatePost(ctx context.Context, request *genproto.CreatePostRequest) (*genproto.Post, error) {
	return p.UnimplementedPostServiceServer.CreatePost(ctx, request)
}

func (p postService) UpdatePost(ctx context.Context, post *genproto.Post) (*genproto.Post, error) {
	return p.UnimplementedPostServiceServer.UpdatePost(ctx, post)
}

func (p postService) DeletePost(ctx context.Context, request *genproto.GetPostRequest) (*genproto.Empty, error) {
	return p.UnimplementedPostServiceServer.DeletePost(ctx, request)
}

func (p postService) FilterPost(ctx context.Context, request *genproto.FilterPostRequest) (*genproto.ListPostsResponse, error) {
	type postResponse struct {
		*genproto.Post
		UserID uint32 `json:"userId"`
	}

	type response []*postResponse

	// build url
	u := url.Values{}
	u.Set("userId", strconv.Itoa(int(request.GetUserId())))

	var results response
	if err := p.fetch(ctx, http.MethodGet, p.serverURL+"/posts?"+u.Encode(), nil, &results); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// fetch users
	var (
		wg     sync.WaitGroup
		client = genproto.NewUserServiceClient(p.userServiceConn)
		r      = &genproto.ListPostsResponse{
			Posts: make([]*genproto.Post, len(results)),
		}
	)

	for i, p := range results {
		wg.Add(1)
		go func(post *postResponse, j int) {
			defer wg.Done()

			var err error
			post.Post.User, err = client.GetUser(ctx, &genproto.GetUserRequest{Id: post.UserID})
			if err != nil {
				log.Errorf("client.GetUser: %v", err)
				return
			}

			r.Posts[j] = post.Post
		}(p, i)
	}

	wg.Wait()
	return r, nil
}
