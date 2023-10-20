package main

import (
	"rest-api/connection"
	"rest-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	connection.InitConnection()
	route := gin.Default()
	routes.SetRoutes(route)
	route.Run("localhost:8080")
}
