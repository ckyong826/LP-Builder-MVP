package routes

import (
	"backend/internal/controllers"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, container *services.ServiceContainer) {
	
    userController := controllers.NewUserController(container.UserService)
    userRoutes := router.Group("/users")
    {
        userRoutes.GET("/", userController.FindAll)
        userRoutes.GET("/:id", userController.FindOneById)
        userRoutes.POST("/", userController.Create)
        userRoutes.PATCH("/:id", userController.Update)
        userRoutes.DELETE("/:id", userController.Delete)
    }

    templateController:=  controllers.NewTemplateController(container.TemplateService)
    templateRoutes := router.Group("/templates")
    {
        templateRoutes.GET("/", templateController.FindAll)
        templateRoutes.GET("/:id", templateController.FindOneById)
        templateRoutes.POST("/", templateController.Create)
        templateRoutes.POST("/convert-url", templateController.ConvertUrlToFile)
        templateRoutes.PATCH("/:id", templateController.Update)
        templateRoutes.DELETE("/:id", templateController.Delete)
    }

}
