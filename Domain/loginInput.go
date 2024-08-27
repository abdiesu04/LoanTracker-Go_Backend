package Domain


type LoginInput struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}