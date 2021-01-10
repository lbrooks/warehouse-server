package item

import (
	"context"
	"main/utils/tracing"

	"github.com/gin-gonic/gin"
)

// AddRoutes Add Testing Routes
func AddRoutes(e *gin.Engine, api *gin.RouterGroup, act *gin.RouterGroup, s Service) {
	addAPIRoutes(api.Group("item"), s)
	addTemplateRoutes(e, e.Group("item"), s)
	addRedirectRoutes(act.Group("item"), s)
}

func addAPIRoutes(rg *gin.RouterGroup, s Service) {
	rg.GET("", tracing.JSONRoute("routes-item", "get-all", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
		return s.GetAllItems(spanCtx)
	}))

	rg.POST("incr", tracing.JSONRoute("routes-item", "increment", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
		barcode := ginCtx.Query("barcode")
		name := ginCtx.Query("name")
		return s.AdjustQuantity(spanCtx, barcode, name, 1)
	}))

	rg.POST("decr", tracing.JSONRoute("routes-item", "decrement", func(spanCtx context.Context, ginCtx *gin.Context) (interface{}, error) {
		barcode := ginCtx.Query("barcode")
		name := ginCtx.Query("name")
		return s.AdjustQuantity(spanCtx, barcode, name, -1)
	}))
}

func addTemplateRoutes(e *gin.Engine, rg *gin.RouterGroup, s Service) {
	e.LoadHTMLGlob("item/templates/*")

	rg.GET("", tracing.HTMLRoute("templates-item", "list", func(spanCtx context.Context, ginCtx *gin.Context) (string, interface{}) {
		items, _ := s.GetAllItems(spanCtx)

		return "list.html", gin.H{
			"title": "Main website",
			"items": items,
		}
	}))

	rg.GET("error", tracing.HTMLRoute("templates-error", "error", func(spanCtx context.Context, ginCtx *gin.Context) (string, interface{}) {
		return "error.html", gin.H{
			"title": "Main website",
		}
	}))
}

func addRedirectRoutes(rg *gin.RouterGroup, s Service) {
	rg.GET("incr", tracing.RedirectRoute("redirect-item", "increment", func(spanCtx context.Context, ginCtx *gin.Context) string {
		barcode := ginCtx.Query("barcode")
		name := ginCtx.Query("name")

		_, err := s.AdjustQuantity(spanCtx, barcode, name, 1)
		if err != nil {
			return "/item/error"
		}

		return "/item"
	}))

	rg.GET("decr", tracing.RedirectRoute("redirect-item", "decrement", func(spanCtx context.Context, ginCtx *gin.Context) string {
		barcode := ginCtx.Query("barcode")
		name := ginCtx.Query("name")

		_, err := s.AdjustQuantity(spanCtx, barcode, name, -1)
		if err != nil {
			return "/item/error"
		}

		return "/item"
	}))
}
