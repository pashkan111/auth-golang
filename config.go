package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Db() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	port, _ := strconv.Atoi(os.Getenv("MONGO_PORT"))
	host := os.Getenv("MONGO_HOST")
	connectionString := fmt.Sprintf("mongodb://%s:%d", host, port)
	fmt.Println(connectionString)
	clientOptions := options.Client().ApplyURI(connectionString)
	client, connect_err := mongo.Connect(context.Background(), clientOptions)
	if connect_err != nil {
		log.Fatal(connect_err)
	}

	ping_err := client.Ping(context.Background(), nil)
	if ping_err != nil {
		log.Fatal(ping_err)
	}
	log.Println("Connected to MongoDB!")
	return client
}
