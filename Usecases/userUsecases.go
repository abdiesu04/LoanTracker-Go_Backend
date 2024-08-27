package Usecases

import (
    "LoanTracker/Domain"
    "LoanTracker/Repository"
    "LoanTracker/infrastructure"
    "errors"
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUsecase interface {
    CreateUser(user *Domain.RegisterInput) (*Domain.User, error)
    Verify(token string) error
    Login(c *gin.Context, LoginUser *Domain.LoginInput) (string, string, error)
    FindUser() ([]Domain.User, error)
    ForgotPassword(c *gin.Context, username string) (string, error)
    Reset(c *gin.Context, token string) (string, error)
    UpdatePassword(username string, newPassword string) error
    GetMyProfile(username string) (*Domain.User, error)
    DeleteUser(id string) error
}

type userUsecase struct {
    userRepository  Repository.UserRepository
    logRepository   Repository.LogRepository
    emailService    *infrastructure.EmailService
    passwordService *infrastructure.PasswordService
}

func NewUserUsecase(userRepo Repository.UserRepository, logRepo Repository.LogRepository, emailService *infrastructure.EmailService) UserUsecase {
    return &userUsecase{
        userRepository:  userRepo,
        logRepository:   logRepo,
        emailService:    emailService,
        passwordService: infrastructure.NewPasswordService(),
    }
}

func (u *userUsecase) CreateUser(user *Domain.RegisterInput) (*Domain.User, error) {
    if existingUser, err := u.userRepository.FindByEmail(user.Email); err != nil {
        return nil, err
    } else if existingUser != nil {
        return nil, errors.New("email already registered")
    }

    if existingUser, err := u.userRepository.FindByUsername(user.Username); err != nil {
        return nil, err
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

    newToken, err := infrastructure.GenerateResetToken(userData.Username, userData.Role)
    if err != nil {
        return nil, fmt.Errorf("failed to generate reset token: %v", err)
    }

    verifyURL := fmt.Sprintf("http://localhost:8080/users/verify-email/%s", newToken)

    emailTitle := "Verify Your Account"
    emailBody := fmt.Sprintf(`<html>
    <!-- Email content here -->
    `, verifyURL, verifyURL)

    err = u.emailService.SendEmail(userData.Email, emailTitle, emailBody)
    if err != nil {
        return nil, fmt.Errorf("failed to send verification email: %v", err)
    }

    createdUser, err := u.userRepository.CreateUser(userData)
    if err != nil {
        return nil, err
    }

    // Log user creation
    log := &Domain.SystemLog{
        Timestamp: time.Now(),
        UserID:    createdUser.ID.Hex(),
        Action:    "User Registered",
        Details:   fmt.Sprintf("User %s registered", userData.Username),
    }
    _ = u.logRepository.CreateLog(log)

    return createdUser, nil
}


func (u *userUsecase) Verify(token string) error {
    claims, err := infrastructure.ParseResetToken(token)
    if err != nil {
        fmt.Println("Error parsing token:", err)
        return err
    }

    user, err := u.userRepository.FindByUsername(claims.Username)
    if err != nil {
        return errors.New("user not found")
    }

    err = u.userRepository.Update(user.Username, bson.M{"is_active": true})
    if err != nil {
        return fmt.Errorf("failed to verify user: %v", err)
    }

    // Log the verification action
    log := &Domain.SystemLog{
        Timestamp: time.Now(),
        UserID:    user.ID.Hex(),
        Action:    "User Verified",
        Details:   fmt.Sprintf("User %s verified their email", user.Username),
    }
    _ = u.logRepository.CreateLog(log)

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

    // Log the login action
    log := &Domain.SystemLog{
        Timestamp: time.Now(),
        UserID:    user.ID.Hex(),
        Action:    "User Logged In",
        Details:   fmt.Sprintf("User %s logged in", user.Username),
    }
    _ = u.logRepository.CreateLog(log)

    return accessToken, refreshToken, nil
}

func (u *userUsecase) FindUser() ([]Domain.User, error) {
	users, err := u.userRepository.ShowUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}


func (u *userUsecase) GetMyProfile(username string) (*Domain.User, error) {
	user, err := u.userRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) ForgotPassword(c *gin.Context, username string) (string, error) {
	user, err := u.userRepository.FindByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	resetToken, err := infrastructure.GenerateResetToken(user.Username, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %v", err)
	}


	subject := "Password Reset Request"
	body := fmt.Sprintf(`
	Hi %s,

	It seems like you requested a password reset. No worries, it happens to the best of us! You can reset your password by clicking the link below:

	<a href="http://localhost:8080/users/reset-password/%s">Reset Your Password</a>

	If you did not request a password reset, please ignore this email.

Best regards,
	Your Support Team
	`, user.FirstName, resetToken)

	err = u.emailService.SendEmail(user.Email, subject, body)
	if err != nil {
		return "", fmt.Errorf("failed to send reset email: %v", err)
	}

	return resetToken, nil
}


func (u *userUsecase) Reset(c *gin.Context, token string) (string, error) {

	claims, err := infrastructure.ParseResetToken(token)
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return "", err
	}

	user, err := u.userRepository.FindByUsername(claims.Username)

	if err != nil {
		return "", errors.New("user not found")
	}

	refreshToken, err := infrastructure.GenerateRefreshToken(user.Username)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	c.SetCookie("refresh_token", refreshToken, 3600, "/", "", false, true)

	access_token, err := infrastructure.GenerateJWT(user.Username, user.Role)

	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %v", err)
	}

	return access_token, nil
}


func (u *userUsecase) UpdatePassword(username string, newPassword string) error {

	hashedPassword, err := u.passwordService.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	err = u.userRepository.Update(username, bson.M{"password": hashedPassword})
	if err != nil {
		// Log the failure to update the password
		log := &Domain.SystemLog{
			Timestamp: time.Now(),
			UserID:   username, 
			Action:    "Update Password",
			Details:   fmt.Sprintf("Failed to update password for user: %s", username),
		}
		u.logRepository.CreateLog(log)

		return fmt.Errorf("failed to update password: %v", err)
	}

	// Log the successful password update
	log := &Domain.SystemLog{
		Timestamp: time.Now(),
		UserID:    username, 
		Action:    "Update Password",
		Details:   fmt.Sprintf("Password updated successfully for user: %s", username),
	}
	u.logRepository.CreateLog(log)

	return nil
}


func (u *userUsecase) DeleteUser(id string) error {
	_, err := u.userRepository.FindByID(id)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	err = u.userRepository.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}