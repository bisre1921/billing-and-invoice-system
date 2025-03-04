package main

import (
	"fmt"
	"log"

	"github.com/bisre1921/billing-and-invoice-system/config"
)

func main() {
	fmt.Println("Starting server...")

	err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	} else {
		fmt.Println("Database connection established successfully!")
	}

}
