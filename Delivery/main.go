package main

import (
	"LoanTracker/Delivery/Controller"
	"LoanTracker/Delivery/router"
	"LoanTracker/Repository"
	"LoanTracker/Usecases"
	"LoanTracker/infrastructure"
	"fmt"

	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)




func main(){
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("++++++++++++++++++++++++++")
			log.Fatal("Error loading .env file")
		}
	} else {
		log.Fatal("No .env file found")
	}
	mongoURI := os.Getenv("MONGO_URL")
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	userDatabase := client.Database("LoanTracker")

	userCollection := userDatabase.Collection("User")
	// blogCollection := userDatabase.Collection("Blog")

	userRepository := Repository.NewUserRepository(userCollection)

	// blogRepository := Repository.NewBlogRepository(blogCollection)
	// blogUsecase := Usecases.NewBlogUsecase(blogRepository)

	// Initialize the Email Service
	emailService := infrastructure.NewEmailService()

	// Initialize the User Usecase with the User Repository and Email Service
	userUsecase := Usecases.NewUserUsecase(userRepository, emailService)
	// blogController := controller.NewBlogController(blogUsecase, userUsecase)
	userController := Controller.NewUserController(userUsecase)

	// Start the token cleanup cron job

	router := router.SetupRouter(userController)
	log.Fatal(router.Run(":8080"))

}
