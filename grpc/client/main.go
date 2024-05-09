package main

import (
	"context"
	"fmt"
	"go-jeager/tracing"
	grpcHeader "go-jeager/tracing/header/grpc"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	ecpb "google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:8000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	fmt.Println("--- calling helloworld.Greeter/SayHello ---")
	hwc := ecpb.NewEchoClient(conn)
	tracingService := tracing.NewTracingService(context.Background(), "client grpc")
	for i := 0; i < 1000; i++ {
		var ctx context.Context
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		ctx, span := tracingService.Start(ctx, "client process")
		subCtx, subSpan := tracingService.Start(ctx, "client get data")
		p := propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
		mdRequest := grpcHeader.Header(metadata.New(nil))
		p.Inject(subCtx, mdRequest)
		subCtx = metadata.NewOutgoingContext(subCtx, metadata.MD(mdRequest))
		r, err := hwc.UnaryEcho(subCtx, &ecpb.EchoRequest{Message: "this is examples/tracing"})
		if err != nil {
			fmt.Println(err)
		}
		if r != nil {
			fmt.Println(r.Message)
		}
		subSpan.End()
		_, subSpan2 := tracingService.Start(ctx, "client internal process")
		min := 10
		max := 50
		timeRandom := rand.Intn(max-min) + min
		time.Sleep(time.Duration(timeRandom) * time.Millisecond)
		subSpan2.End()
		span.End()
	}
}
