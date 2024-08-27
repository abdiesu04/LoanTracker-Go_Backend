package router

import (
	"LoanTracker/Delivery/Controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userController *Controller.UserController) *gin.Engine {
	router := gin.Default()
	userRouter := router.Group("/users")
	userRouter.POST("/register", userController.Register)
	userRouter.GET("/verify-email/:token", userController.VerifyEmail)
	return router
}
