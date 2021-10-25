package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/SkYNewZ/grpc-jsonplaceholder/internal/genproto"

	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const listenPort = "5051"

func init() {
	log.SetLevel(log.TraceLevel)
	if os.Getenv("K_SERVICE") != "" {
		log.SetFormatter(stackdriver.NewFormatter())
	}
}

func main() {
	port := listenPort
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	var (
		svc   = newPostService()
		entry = log.NewEntry(log.StandardLogger())
	)

	// define some logrus options
	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
		grpc_logrus.WithDurationField(grpc_logrus.DefaultDurationToField),
		grpc_logrus.WithCodes(grpc_logging.DefaultErrorToCode),
	}

	log.Debugln("initializing gRPC server")
	var srv = grpc.NewServer(
		grpc.UnaryInterceptor(middleware.ChainUnaryServer(grpc_logrus.UnaryServerInterceptor(entry, opts...))),
		grpc.StreamInterceptor(middleware.ChainStreamServer(grpc_logrus.StreamServerInterceptor(entry, opts...))),
	)

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalln(err)
	}

	genproto.RegisterPostServiceServer(srv, svc)
	healthpb.RegisterHealthServer(srv, svc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("starting to listen on tcp: %q", listen.Addr().String())
		if err := srv.Serve(listen); err != nil {
			log.Fatalln(err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
	stop()

	srv.GracefulStop()
	log.Println("server exiting")
}

func mustMapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v == "" {
		log.Fatalf("environment variable %q not set", envKey)
	}
	*target = v
}

func mustConnGRPC(ctx context.Context, conn **grpc.ClientConn, addr string, opts ...grpc.DialOption) {
	var err error
	*conn, err = grpc.DialContext(ctx, addr, opts...)

	if err != nil {
		log.WithError(err).Fatalf("grpc: failed to connect %s", addr)
	}
}
