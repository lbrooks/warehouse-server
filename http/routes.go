package http

import (
	"context"
	"encoding/json"

	"main/inventory"
	"main/inventory/tracing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/label"
)

const traceName = "item.routes"

// AddRoutes Add Testing Routes
func AddRoutes(api *gin.RouterGroup, s inventory.ItemService) {
	itemAPIRoutes := api.Group("item")

	itemAPIRoutes.GET("", tracing.JSONRoute("routes-item", "get", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
		span := tracing.GetSpan(spanCtx)

		var item inventory.Item
		err := ginCtx.Bind(&item)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}

		{
			reqJSON, _ := json.Marshal(item)
			span.SetAttributes(
				label.String("request.body", string(reqJSON)),
			)
		}

		return s.Search(spanCtx, item)
	}))

	itemAPIRoutes.POST("update", tracing.JSONRoute("routes-item", "update", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
		span := tracing.GetSpan(spanCtx)

		var item inventory.Item
		err := ginCtx.Bind(&item)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}

		{
			reqJSON, _ := json.Marshal(item)
			span.SetAttributes(
				label.String("request.body", string(reqJSON)),
			)
		}

		return s.Update(spanCtx, item)
	}))
}
