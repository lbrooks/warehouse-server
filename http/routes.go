package http

import (
	"context"
	"encoding/json"

	"github.com/lbrooks/warehouse"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/label"
)

// AddRoutes Add Testing Routes
func AddRoutes(api *gin.RouterGroup, s warehouse.ItemService) {
	itemAPIRoutes := api.Group("item")

	itemAPIRoutes.GET("", JSONRoute("routes", "get", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
		span := warehouse.GetSpan(spanCtx)

		var item warehouse.Item
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

	itemAPIRoutes.POST("update", JSONRoute("routes", "update", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
		span := warehouse.GetSpan(spanCtx)

		var item warehouse.Item
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
