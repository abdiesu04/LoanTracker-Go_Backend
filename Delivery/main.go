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
	loanCollection := userDatabase.Collection("Loan")


	userRepository := Repository.NewUserRepository(userCollection)

	loanRepository := Repository.NewLoanRepository(loanCollection)
	loanUsecase := Usecases.NewLoanUsecase(loanRepository)

	// Initialize the Email Service
	emailService := infrastructure.NewEmailService()

	userUsecase := Usecases.NewUserUsecase(userRepository, emailService)
	loanController := Controller.NewLoanController(loanUsecase)
	userController := Controller.NewUserController(userUsecase)


	router := router.SetupRouter(userController , loanController)
	log.Fatal(router.Run(":8080"))

}
