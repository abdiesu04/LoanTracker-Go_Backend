package Domain

type RegisterInput struct {
	FirstName string `json:"first_name" bson:"first_name"`
	LastName string `json:"last_name" bson:"last_name"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Email string `json:"email" bson:"email"`
	
}