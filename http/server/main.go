package main

import (
	"context"
	"fmt"
	"go-jeager/tracing"
	httpHeader "go-jeager/tracing/header/http"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/propagation"
)

func startHttpServer(port int) {
	tracingService := tracing.NewTracingService(context.Background(), "server http")
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		p := propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
		ctx := p.Extract(r.Context(), httpHeader.Header(r.Header))
		_, span := tracingService.Start(ctx, "server response")
		defer span.End()
		min := 200
		max := 300
		timeRandom := rand.Intn(max-min) + min
		time.Sleep(time.Duration(timeRandom) * time.Millisecond)
		fmt.Fprintf(w, fmt.Sprintf("response from %s:%d", "localhost", port))
	})
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 3 * time.Second, //nolint:gomnd // common
	}
	fmt.Printf("Start http server listening:%d", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}

func main() {
	startHttpServer(8000)
}
