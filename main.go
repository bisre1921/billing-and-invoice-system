package main

import (
	"fmt"
	"log"

	"github.com/bisre1921/billing-and-invoice-system/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting server...")

	err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	} else {
		fmt.Println("Database connection established successfully!")
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}
