package main

import (
	"context"
	"log"
	"os"

	"github.com/lbrooks/warehouse"
	"github.com/lbrooks/warehouse/http"
	"github.com/lbrooks/warehouse/memory"
	"github.com/lbrooks/warehouse/server"

	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	flush := warehouse.InitializeJaeger()
	defer flush()

	webServer := server.NewWebServer()
	apiRoutes := webServer.Group("api")

	itemService := memory.NewItemService(context.Background(), true)
	http.AddRoutes(apiRoutes, itemService)

	err := webServer.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		os.Exit(1)
	}
}
