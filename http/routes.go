package http

import (
	"context"

	"github.com/lbrooks/warehouse"

	"github.com/gin-gonic/gin"
)

// AddRoutes Add Testing Routes
func AddRoutes(api *gin.RouterGroup, s warehouse.ItemService) {
	itemAPIRoutes := api.Group("item")

	itemAPIRoutes.GET("", TracedRoute("routes", "get", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
		var item warehouse.Item
		err := ginCtx.Bind(&item)
		if err != nil {
			return nil, err
		}

		return s.Search(spanCtx, item)
	}))

	itemAPIRoutes.POST("update", TracedRoute("routes", "update", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
		var item warehouse.Item
		err := ginCtx.Bind(&item)
		if err != nil {
			return nil, err
		}

		return s.Update(spanCtx, item)
	}))
}
