package internal

import (
	"backend/config"
	"backend/internal/middleware"
	"backend/internal/routes"
	"backend/internal/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	gin.ForceConsoleColor()

	cfg, err := config.LoadConfig()
	if err != nil {
        log.Fatal("Failed to load configuration:", err)
    }

    db, err := config.InitDB(cfg)
    if err != nil {
        log.Fatal("Failed to initialize database:", err)
    }

    serviceContainer := services.NewServiceContainer(db)

    router := gin.New() 
    router.Use(gin.Recovery())  
    router.Use(middleware.RequestLogger()) 
    
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