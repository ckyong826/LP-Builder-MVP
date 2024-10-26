package main

import (
	"backend/src/api"
	"backend/src/config"
)

func main() {
    config.LoadConfig() // load environment variables

    router := api.InitRouter() // initialize routes
    router.Run(":8080")        // start the server
}
