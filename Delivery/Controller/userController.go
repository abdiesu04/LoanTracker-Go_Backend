package Controller

import (
	"LoanTracker/Domain"
	"LoanTracker/Usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase Usecases.UserUsecase
}

func NewUserController(userUsecase Usecases.UserUsecase) *UserController {
	return &UserController{
		UserUsecase: userUsecase,
	}
}

func (u *UserController) Register(c *gin.Context) {
	var input Domain.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := u.UserUsecase.CreateUser(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}


func (u *UserController) VerifyEmail(c *gin.Context) {
	token := c.Param("token")
	err := u.UserUsecase.Verify(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified"})
}