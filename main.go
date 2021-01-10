package main

import (
	"context"
	"main/item"
	"main/utils/tracing"

	"github.com/gin-gonic/gin"
)

func main() {
	flush := tracing.Initialize()
	defer flush()

	r := gin.Default()

	apiRoutes := r.Group("api")
	actionRoutes := r.Group("action")

	itemService := item.NewService(item.NewDaoInMemory(context.Background(), true))
	item.AddRoutes(r, apiRoutes, actionRoutes, itemService)
	r.Run() // listen and serve on 0.0.0.0:8080
}
