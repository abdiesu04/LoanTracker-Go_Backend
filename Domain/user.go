package Domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type  User struct {
	ID  primitive.ObjectID `json:"id" bson:"_id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName string `json:"last_name" bson:"last_name"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Email string `json:"email" bson:"email"`
	IsActive bool `json:"is_active" bson:"is_active"`
	Role string `json:"role" bson:"role"`

}

type ForgetPasswordInput struct {
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
}

type ResetPasswordInput struct {
	Username    string `json:"username" bson:"username"`
	NewPassword string `json:"password" bson:"password"`
}

type ChangePasswordInput struct {
	NewPassword string `json:"password" bson:"password"`
}