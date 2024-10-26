package routes

import (
	"backend/src/api/controllers"
	"backend/src/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
    userRoutes := router.Group("/users", middleware.AuthMiddleware())
    {
        userRoutes.GET("/", controllers.GetUsers)
        userRoutes.POST("/", controllers.CreateUser)
    }
}
