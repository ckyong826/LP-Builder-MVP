package internal

import (
	"backend/config"
	"backend/internal/routes"
	"backend/internal/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func StartServer() {
    config.LoadConfig()
    config.InitDB()

    serviceContainer := services.NewServiceContainer(config.DB)
    router := gin.Default()

    routes.RegisterRoutes(router, serviceContainer)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Starting server on port %s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
