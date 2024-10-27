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

}
