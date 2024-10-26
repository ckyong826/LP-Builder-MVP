package main

import (
	"backend/src/api"
	"backend/src/config"
	"backend/src/database"
	"backend/src/database/migrations"
)

func main() {
    config.LoadConfig() 
    database.Connect()
    migrations.Migrate()
    router := api.InitRouter() 
    router.Run(":8080")        
}
