package main

import (
	"LoanTracker/Delivery/Controller"
	"LoanTracker/Delivery/router"
	"LoanTracker/Repository"
	"LoanTracker/Usecases"
	"LoanTracker/infrastructure"

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
	logCollection := userDatabase.Collection("Logs")


	userRepository := Repository.NewUserRepository(userCollection )

	loanRepository := Repository.NewLoanRepository(loanCollection)
	logRepository := Repository.NewLogRepository(logCollection)
	loanUsecase := Usecases.NewLoanUsecase(loanRepository , logRepository)

	// Initialize the Email Service
	emailService := infrastructure.NewEmailService()

	userUsecase := Usecases.NewUserUsecase(userRepository, logRepository , emailService ) 
	loanController := Controller.NewLoanController(loanUsecase)
	userController := Controller.NewUserController(userUsecase)
	logController := Controller.NewLogController(logRepository)


	router := router.SetupRouter(userController , loanController , logController)
	log.Fatal(router.Run(":8080"))

}
