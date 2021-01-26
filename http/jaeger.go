package http

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lbrooks/warehouse"
	"go.opentelemetry.io/otel/label"
)

// TracedRouteHandler function to handle request
type TracedRouteHandler func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error)

// JSONRoute wrapper to trace route
func JSONRoute(tracerName, operationName string, handler TracedRouteHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		sc, span := warehouse.CreateSpan(c, tracerName, operationName)
		defer span.End()

		for k, v := range c.Request.Header {
			span.SetAttributes(label.String("request.header."+k, strings.Join(v, " | ")))
		}
		span.SetAttributes(label.String("request.query", c.Request.URL.RawQuery))

		msg, err := handler(sc, c)

		if err != nil {
			span.SetAttributes(label.Int("http.status_code", 500))
			span.RecordError(err)

			c.JSON(500, err.Error())
			return
		}

		span.SetAttributes(label.Int("http.status_code", 200))
		c.JSON(200, msg)
	}
}
