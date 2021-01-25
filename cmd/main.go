package cmd

import (
	"context"
	"os"
	"time"

	"github.com/lbrooks/inventory/server/dao"
	"github.com/lbrooks/inventory/server/tracing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	flush := tracing.Initialize()
	defer flush()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "OPTIONS", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	apiRoutes := r.Group("api")

	itemService := item.NewService(dao.NewDaoInMemory(context.Background(), true))
	item.AddRoutes(apiRoutes, itemService)

	err := r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		os.Exit(1)
	}
}
