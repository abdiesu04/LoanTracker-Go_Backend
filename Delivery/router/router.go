package router

import (
	"LoanTracker/Delivery/Controller"
	"LoanTracker/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userController *Controller.UserController , loanController *Controller.LoanController) *gin.Engine {
	router := gin.Default()
    
	loanRouter := router.Group("/").Use(infrastructure.AuthMiddleware())
	loanRouter.POST("/loans", loanController.ApplyLoanHandler)
	loanRouter.GET("/loans/:id", loanController.GetLoanStatusHandler)
	loanRouter.GET("/loans", loanController.GetAllLoansHandler)
	loanRouter.PATCH("/loans/:id/:status", loanController.UpdateLoanStatusHandler)
	loanRouter.DELETE("/loans/:id", loanController.DeleteLoanHandler)


	userRouter := router.Group("/users")
	userRouter.POST("/register", userController.Register)
	userRouter.GET("/verify-email/:token", userController.VerifyEmail)
	userRouter.POST("/login", userController.Login)
	userRouter.POST("/refresh-token", userController.RefreshToken)
	userRouter.GET("/profile", userController.GetProfile)
	userRouter.POST("/password-reset", userController.ForgotPassword)
	userRouter.GET("/reset-password/:token", userController.ResetPassword)


	authRoutes := userRouter.Use(infrastructure.AuthMiddleware())	
	authRoutes.POST("/password-update", userController.ChangePassword)
	


	adminRoutes := router.Group("/admin")
	adminRoutes.Use(infrastructure.RoleMiddleware("admin"))
	adminRoutes.GET("/users", userController.FindUsers)
	// r.GET("/loans/:id", loanController.GetLoanStatusHandler)
	// r.GET("/admin/loans", loanController.GetAllLoansHandler)
	// r.PATCH("/admin/loans/:id/status", loanController.UpdateLoanStatusHandler)
	// r.DELETE("/admin/loans/:id", loanController.DeleteLoanHandler)

	return router
}
