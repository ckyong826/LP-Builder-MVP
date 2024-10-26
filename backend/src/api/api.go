package api

import (
	"backend/src/api/routes"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
    router := gin.Default()
    routes.InitRoutes(router)
    return router
}
