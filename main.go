package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang-showcase-qr-backend/controller"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Get MongoDB connection string, database and collection name from environment variables
	uri := os.Getenv("MONGODB_URI")
	databaseName := os.Getenv("MONGODB_DATABASE")
	collectionName := os.Getenv("COLLECTION_NAME")

	// Set up a MongoDB connection
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	// Create router and attach routes
	r := mux.NewRouter()
	userRouter := r.NewRoute().Subrouter()
	user.NewUserController(userRouter, client.Database(databaseName).Collection(collectionName))

	// Set the NotFoundHandler for unmatched requests
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Unmatched request: %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	})

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8080", r))
}
