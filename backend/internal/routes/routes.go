package routes

import (
	"backend/internal/controllers"
	"backend/internal/middleware"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, container *services.ServiceContainer) {
    // CORS middleware
    router.Use(middleware.CORS())

    // API version group
    api := router.Group("/api")
    
    // User routes
    userController := controllers.NewUserController(container.UserService)
    users := api.Group("/users")
    {
        users.GET("", userController.FindAll)
        users.GET("/:id", userController.FindOneById)
        users.POST("", userController.Create)
        users.PUT("/:id", userController.Update)  // Changed from PATCH to PUT to match controller
        users.DELETE("/:id", userController.Delete)
    }

    // Template routes
    templateController := controllers.NewTemplateController(container.TemplateService)
    templates := api.Group("/templates")
    {
        templates.GET("", templateController.FindAll)
        templates.GET("/:id", templateController.FindOneById)
        templates.GET("/:id/content", templateController.GetTemplateContent)
        templates.POST("", templateController.Create)
        templates.POST("/convert", templateController.ConvertUrlToFile)  // Changed URL to match controller
        templates.PUT("/:id", templateController.Update)  // Changed from PATCH to PUT to match controller
        templates.DELETE("/:id", templateController.Delete)
    }
}