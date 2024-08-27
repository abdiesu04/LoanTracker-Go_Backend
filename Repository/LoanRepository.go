package Repository

import (
	"LoanTracker/Domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoanRepository interface {
	ApplyLoan(loan *Domain.Loan) (*Domain.ApplyLoanResponse, error)
	GetLoanByID(id primitive.ObjectID) (*Domain.Loan, error)
	GetAllLoans(status string, order string) ([]Domain.Loan, error)
	UpdateLoanStatus(id primitive.ObjectID, status string) error
	DeleteLoan(id primitive.ObjectID) error
}

type loanRepository struct {
	collection *mongo.Collection
}

func NewLoanRepository(collection *mongo.Collection) LoanRepository {
	return &loanRepository{collection}
}

func (r *loanRepository) ApplyLoan(loan *Domain.Loan) (*Domain.ApplyLoanResponse, error) {
	loan.ID = primitive.NewObjectID()
	loan.Status = "pending"
	loan.CreatedAt = time.Now().Unix()
	loan.UpdatedAt = time.Now().Unix()

	_, err := r.collection.InsertOne(context.TODO(), loan)
	if err != nil {
		return nil, err
	}

	response := &Domain.ApplyLoanResponse{
		LoanID:  loan.ID.Hex(),
		Status:  loan.Status,
		Message: "Your loan application has been submitted successfully.",
	}
	return response, nil
}

func (r *loanRepository) GetLoanByID(id primitive.ObjectID) (*Domain.Loan, error) {
	var loan Domain.Loan
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&loan)
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *loanRepository) GetAllLoans(status string, order string) ([]Domain.Loan, error) {
	var loans []Domain.Loan
	filter := bson.M{}
	if status != "" && status != "all" {
		filter["status"] = status
	}

	sortOrder := 1
	if order == "desc" {
		sortOrder = -1
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", sortOrder}})

	cursor, err := r.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &loans); err != nil {
		return nil, err
	}
	return loans, nil
}
func (r *loanRepository) UpdateLoanStatus(id primitive.ObjectID, status string) error {
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updated_at": time.Now().Unix(),
		},
	}
	_, err := r.collection.UpdateByID(context.TODO(), id, update)
	return err
}

func (r *loanRepository) DeleteLoan(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}
