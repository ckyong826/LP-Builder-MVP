package api

import (
	"backend/src/api/controllers"
	"backend/src/api/routes"
	"backend/src/api/services"
	"backend/src/database"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	db := database.DB 

	// Users
	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)

	// Initialize routes
	routes.InitRoutes(router, userController)

	return router
}
