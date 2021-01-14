package item

import (
	"context"
	"encoding/json"
	"main/utils/tracing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/label"
)

// AddRoutes Add Testing Routes
func AddRoutes(api *gin.RouterGroup, s Service) {
	addAPIRoutes(api.Group("item"), s)
}

func addAPIRoutes(rg *gin.RouterGroup, s Service) {
	rg.GET("", tracing.JSONRoute("routes-item", "get-all", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
		return s.GetAllItems(spanCtx)
	}))

	rg.POST("update", tracing.JSONRoute("routes-item", "update", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
		span := tracing.GetSpan(spanCtx)

		var item Item
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
