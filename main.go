package main

import (
	"log"
	"os"
	"rest-api/connection"
	"rest-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connection.InitConnection()
	route := gin.Default()
	routes.SetRoutes(route)
	route.Run(os.Getenv("PORT"))
}
