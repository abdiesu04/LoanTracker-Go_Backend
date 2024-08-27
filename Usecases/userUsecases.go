package Usecases

import (
	"LoanTracker/Domain"
	"LoanTracker/Repository"
	"LoanTracker/infrastructure"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type UserUsecase interface {
	CreateUser(user *Domain.RegisterInput) (*Domain.User, error)
	Verify(token string) error
}


type userUsecase struct {
	userRepository  Repository.UserRepository
	emailService    *infrastructure.EmailService
	passwordService *infrastructure.PasswordService
}

func NewUserUsecase(userRepo Repository.UserRepository, emailService *infrastructure.EmailService) UserUsecase {
	return &userUsecase{
		userRepository:  userRepo,
		emailService:    emailService,
		passwordService: infrastructure.NewPasswordService(),
	}
}

func (u *userUsecase) CreateUser(user *Domain.RegisterInput) (*Domain.User, error) {

	if existingUser, err := u.userRepository.FindByEmail(user.Email); err != nil {
		return nil, err // Return the actual error
	} else if existingUser != nil {
		return nil, errors.New("email already registered")
	}
	
	if existingUser, err := u.userRepository.FindByUsername(user.Username); err != nil {
		return nil, err // Return the actual error
	} else if existingUser != nil {
		return nil, errors.New("username already registered")
	}
	hashedPassword, err := u.passwordService.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	

	userData := &Domain.User{
		ID:         primitive.NewObjectID(),
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Username:   user.Username,
		Password:   hashedPassword,
		Email:      user.Email,
		IsActive:   false,
	}

	if ok, err := u.userRepository.IsDbEmpty(); ok && err == nil {
		userData.Role = "admin"
	} else {
		userData.Role = "user"
	}
	
	newToken , err := infrastructure.GenerateResetToken(userData.Username , userData.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate reset token: %v", err)
	}

	verifyURL := fmt.Sprintf("http://localhost:8080/users/verify-email/%s", newToken)

	emailTitle := "Verify Your Account"
	emailBody := fmt.Sprintf(`<html>
	<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px;">
		<div style="max-width: 600px; margin: auto; background-color: #ffffff; padding: 20px; border-radius: 10px; box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);">
			<h1 style="color: #007bff; text-align: center;">Welcome to LoanTracker</h1>
			<p style="font-size: 16px; color: #333;">Thank you for registering an account with us. To complete your registration, please verify your account by clicking the button below:</p>
			<div style="text-align: center; margin: 30px 0;">
				<a href="%s" style="display: inline-block; background-color: #007bff; color: #ffffff; padding: 15px 25px; font-size: 16px; font-weight: bold; text-decoration: none; border-radius: 5px;">Click Here to Verify Your Account</a>
			</div>
			<p style="font-size: 14px; color: #666; text-align: center;">If the button above doesn't work, please copy and paste the following URL into your browser:</p>
			<p style="font-size: 14px; color: #007bff; text-align: center;">%s</p>
		</div>
	</body>
	</html>`, verifyURL, verifyURL)


	err = u.emailService.SendEmail(userData.Email, emailTitle, emailBody)
	if err != nil {
		return nil, fmt.Errorf("failed to send verification email: %v", err)
	}




	return u.userRepository.CreateUser(userData)

}

func (u *userUsecase) Verify(token string) error {
	claims, err := infrastructure.ParseResetToken(token)
	if err != nil {
		fmt.Println("Error parsing token:", err)
	}

	user, err := u.userRepository.FindByUsername(claims.Username)
	if err != nil {
		return errors.New("user not found")
	}
	err = u.userRepository.Update(user.Username, bson.M{"is_active": true})
	if err != nil {
		return fmt.Errorf("failed to verify user: %v", err)
	}
	return nil
}


func (u *userUsecase) Login(c *gin.Context, LoginUser *Domain.LoginInput) (string, string, error) {
	user, err := u.userRepository.FindByUsername(LoginUser.Username)
	if err != nil {
		return "", "", errors.New("invalid username or password")
	}

	err = u.passwordService.ComparePasswords(user.Password, LoginUser.Password)
	if err != nil {
		return "", "", errors.New("invalid username or password")
	}

	accessToken, err := infrastructure.GenerateJWT(user.Username, user.Role)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, err := infrastructure.GenerateRefreshToken(user.Username)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	c.SetCookie("refresh_token", refreshToken, 3600, "/", "", false, true)


	if !user.IsActive {
		return "", "", fmt.Errorf("user not verified")
	}

	return accessToken, refreshToken, nil
}