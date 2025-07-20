package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/jyotishmoy12/mongo-golang/controllers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	err :=godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := mongoConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	dbName :=os.Getenv("DB_NAME")
	
	collection := client.Database(dbName).Collection("users")

	
	uc := controllers.NewUserController(collection)

	router := httprouter.New()
	router.GET("/user/:id", uc.GetUser)
	router.POST("/user", uc.CreateUser)
	router.DELETE("/user/:id", uc.DeleteUser)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func mongoConnect() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping to ensure connection works
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	fmt.Println("✅ Connected to MongoDB")
	return client, nil
}
