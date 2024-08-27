package Controller

import (
	"LoanTracker/Domain"
	"LoanTracker/Usecases"
	"LoanTracker/infrastructure"
	"net/http"

	"github.com/dgrijalva/jwt-go"
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

	_, err := u.UserUsecase.CreateUser(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user successfully registered"})
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
func (uc *UserController) Login(c *gin.Context) {
	var input Domain.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	accessToken, refresh_token, err := uc.UserUsecase.Login(c, &input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("refresh_token", refresh_token, 60*60*24*7, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func (uc *UserController) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token not found"})
		return
	}
	var jwtKey = []byte("LoanTracker")

	token, err := jwt.ParseWithClaims(refreshToken, &infrastructure.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return
	}

	// Get the username from the token
	username, err := infrastructure.GetUsernameFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get username from token"})
		c.Abort()
		return
	}

	// Set token claims in context
	claims, ok := token.Claims.(*infrastructure.Claims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
	claims.Username = username

	accessToken, err := infrastructure.GenerateJWT(claims.Username, claims.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	refreshToken, err = infrastructure.GenerateRefreshToken(claims.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func (uc *UserController) GetProfile(c *gin.Context) {
	username := c.GetString("username")
	user, err := uc.UserUsecase.GetMyProfile(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (uc *UserController) FindUsers(c *gin.Context) {

	users, err := uc.UserUsecase.FindUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"users": users})

}
func (uc *UserController) ForgotPassword(c *gin.Context) {

	var input Domain.ForgetPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err := uc.UserUsecase.ForgotPassword(c, input.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset link sent to your email"})

}

func (uc *UserController) ResetPassword(c *gin.Context) {
	reset_token := c.Param("token")

	new_token, err := uc.UserUsecase.Reset(c, reset_token)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": new_token})

}

func (uc *UserController) DeleteUser(c *gin.Context) {
	// username := c.Param("username")
	id := c.GetString("id")

	err := uc.UserUsecase.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (uc *UserController) ChangePassword(c *gin.Context) {
	var input Domain.ChangePasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	username := c.GetString("username")

	err := uc.UserUsecase.UpdatePassword(username, input.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
