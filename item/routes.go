package item

import (
	"context"
	"encoding/json"
	"main/utils/tracing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/label"
)

type adjustJSON struct {
	Barcode string `form:"barcode"`
	Name    string `form:"name" binding:"required"`
}

// AddRoutes Add Testing Routes
func AddRoutes(api *gin.RouterGroup, s Service) {
	addAPIRoutes(api.Group("item"), s)
}

func addAPIRoutes(rg *gin.RouterGroup, s Service) {
	rg.GET("", tracing.JSONRoute("routes-item", "get-all", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
		return s.GetAllItems(spanCtx)
	}))

	rg.POST("incr", tracing.JSONRoute("routes-item", "increment", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
		span := tracing.GetSpan(spanCtx)

		var reqBody adjustJSON
		err := ginCtx.Bind(&reqBody)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}

		{
			reqJSON, _ := json.Marshal(reqBody)
			span.SetAttributes(
				label.String("request.body", string(reqJSON)),
			)
		}

		return s.AdjustQuantity(spanCtx, reqBody.Barcode, reqBody.Name, 1)
	}))

	rg.POST("decr", tracing.JSONRoute("routes-item", "decrement", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
		span := tracing.GetSpan(spanCtx)

		var reqBody adjustJSON
		err := ginCtx.Bind(&reqBody)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}

		{
			reqJSON, _ := json.Marshal(reqBody)
			span.SetAttributes(
				label.String("request.body", string(reqJSON)),
			)
		}

		return s.AdjustQuantity(spanCtx, reqBody.Barcode, reqBody.Name, -1)
	}))
}
