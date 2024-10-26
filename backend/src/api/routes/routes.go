package routes

import (
	"backend/src/api/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine, userController *controllers.UserController) {
	userRoutes := router.Group("/users")
	{
			userRoutes.GET("/", userController.FindAll)       
			userRoutes.POST("/", userController.Create)    
			userRoutes.GET("/:id", userController.FindOneById)     
			userRoutes.PUT("/:id", userController.Update) 
			userRoutes.DELETE("/:id", userController.Delete) 
	}
}
