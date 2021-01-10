package tracing

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const jaegerEndpoint = "http://localhost:14268/api/traces"
const jaegerService = "games-server"
const jaegerHost = "chronos"

// Initialize creates a new trace provider instance and registers it as global trace provider.
func Initialize() func() {
	// Create and install Jaeger export pipeline.
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(jaegerEndpoint),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: jaegerService,
			Tags: []label.KeyValue{
				label.String("host", jaegerHost),
			},
		}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Fatal(err)
	}
	return flush
}

func logRequestsToSpan(span oteltrace.Span, request *http.Request) {
	for k, v := range request.Header {
		span.SetAttributes(label.String("request.header."+k, strings.Join(v, " | ")))
	}
	span.SetAttributes(label.String("request.query", request.URL.RawQuery))
}

// JSONRoute wrapper to trace route
func JSONRoute(tracerName, operationName string, handler func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		sc, span := CreateSpan(c, tracerName, operationName)
		defer span.End()
		logRequestsToSpan(span, c.Request)

		msg, err := handler(sc, c)

		if err != nil {
			span.SetAttributes(label.Int("http.status_code", 500))

			c.JSON(500, err.Error())
			return
		}

		span.SetAttributes(label.Int("http.status_code", 200))
		c.JSON(200, msg)
	}
}

// HTMLRoute wrapper to trace route
func HTMLRoute(tracerName, operationName string, handler func(spanCtx context.Context, ginCtx *gin.Context) (string, interface{})) func(*gin.Context) {
	return func(c *gin.Context) {
		sc, span := CreateSpan(c, tracerName, operationName)
		defer span.End()
		logRequestsToSpan(span, c.Request)

		template, data := handler(sc, c)

		c.HTML(http.StatusOK, template, data)
	}
}

// RedirectRoute wrapper to trace route
func RedirectRoute(tracerName, operationName string, handler func(spanCtx context.Context, ginCtx *gin.Context) string) func(*gin.Context) {
	return func(c *gin.Context) {
		sc, span := CreateSpan(c, tracerName, operationName)
		defer span.End()
		logRequestsToSpan(span, c.Request)

		url := handler(sc, c)

		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// CreateSpan Create Tracing Span from context
func CreateSpan(ctx context.Context, tracerName, operationName string) (context.Context, oteltrace.Span) {
	tr := otel.Tracer(tracerName)
	return tr.Start(ctx, fmt.Sprintf("%s-%s", tracerName, operationName))
}
