package Domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Loan struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Amount    float64            `bson:"amount" json:"amount"`
	Term      int                `bson:"term" json:"term"`
	Purpose   string             `bson:"purpose" json:"purpose"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt int64              `bson:"created_at" json:"created_at"`
	UpdatedAt int64              `bson:"updated_at" json:"updated_at"`
}

type ApplyLoanRequest struct {
	Amount  float64 `json:"amount" binding:"required"`
	Term    int     `json:"term" binding:"required"`
	Purpose string  `json:"purpose" binding:"required"`
}

type ApplyLoanResponse struct {
	LoanID  string `json:"loan_id"`
	Status  string `json:"status"`  
	Message string `json:"message"` 
}
