package config

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Database

func ConnectDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Use local MongoDB URI directly
	localURI := "mongodb://localhost:27017/billing_invoice"
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Connecting to local MongoDB at localhost:27017/billing_invoice...")
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(localURI))
	if err != nil {
		log.Printf("⚠️  Failed to connect to local MongoDB: %v", err)
		log.Println("🚀 Starting server without database connection (development mode)")
		log.Println("   📝 Note: Database operations will not work until connection is established")
		return nil
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("⚠️  Failed to ping local MongoDB: %v", err)
		log.Println("🚀 Starting server without database connection (development mode)")
		log.Println("   📝 Note: Database operations will not work until connection is established")
		return nil
	}

	DB = client.Database("billing_invoice")
	log.Println("✅ Successfully connected to local MongoDB (billing_invoice database)")
	return nil
}
