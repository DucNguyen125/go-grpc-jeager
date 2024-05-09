package main

import (
	"context"
	"fmt"
	"go-jeager/tracing"
	httpHeader "go-jeager/tracing/header/http"
	"io"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/propagation"
)

func main() {
	fmt.Println("--- calling http echo ---")
	tracingService := tracing.NewTracingService(context.Background(), "client http")
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
		req, err := http.NewRequestWithContext(subCtx, http.MethodGet, "http://localhost:8000/echo", nil)
		if err != nil {
			fmt.Println(err)
		}
		p.Inject(subCtx, httpHeader.Header(req.Header))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		if resp != nil {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(body))
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
