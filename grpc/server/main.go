package main

import (
	"context"
	"fmt"
	"go-jeager/tracing"
	grpcHeader "go-jeager/tracing/header/grpc"
	"log"
	"math/rand"
	"net"
	"time"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/metadata"
)

type ecServer struct {
	pb.UnimplementedEchoServer
	addr           string
	tracingService trace.Tracer
}

func (s *ecServer) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	p := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	md, _ := metadata.FromIncomingContext(ctx)
	ctx = p.Extract(ctx, grpcHeader.Header(md))
	_, span := s.tracingService.Start(ctx, "server response")
	defer span.End()
	min := 200
	max := 300
	timeRandom := rand.Intn(max-min) + min
	time.Sleep(time.Duration(timeRandom) * time.Millisecond)
	return &pb.EchoResponse{Message: fmt.Sprintf("%s (from %s)", req.Message, s.addr)}, nil
}

func startGrpcServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	tracingService := tracing.NewTracingService(context.Background(), "server grpc")
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &ecServer{addr: addr, tracingService: tracingService})
	log.Printf("serving on %s\n", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	startGrpcServer(fmt.Sprintf(":%d", 8000))
}
