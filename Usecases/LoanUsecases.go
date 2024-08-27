package Usecases

import (
	"LoanTracker/Domain"
	"LoanTracker/Repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoanUsecase interface {
	ApplyLoan(loanRequest *Domain.ApplyLoanRequest, userID string) (*Domain.ApplyLoanResponse, error)
	GetLoanStatus(id string) (*Domain.Loan, error)
	GetAllLoans(status string, order string) ([]Domain.Loan, error)
	ApproveRejectLoan(id string, status string) error
	DeleteLoan(id string) error
}

type loanUsecase struct {
	loanRepository Repository.LoanRepository
	logRepository  Repository.LogRepository
}

func NewLoanUsecase(loanRepo Repository.LoanRepository, logRepo Repository.LogRepository) LoanUsecase {
	return &loanUsecase{
		loanRepository: loanRepo,
		logRepository:  logRepo,
	}
}

func (u *loanUsecase) ApplyLoan(loanRequest *Domain.ApplyLoanRequest, userID string) (*Domain.ApplyLoanResponse, error) {
	loan := &Domain.Loan{
		UserID:  userID,
		Amount:  loanRequest.Amount,
		Term:    loanRequest.Term,
		Purpose: loanRequest.Purpose,
	}

	response, err := u.loanRepository.ApplyLoan(loan)
	if err != nil {
		return nil, err
	}

	// Log the loan application submission
	log := &Domain.SystemLog{
		Timestamp: time.Now(),
		UserID:    userID,
		Action:    "Loan Application Submission",
		Details:   "User applied for a loan",
	}
	_ = u.logRepository.CreateLog(log)

	return response, nil
}

func (u *loanUsecase) GetLoanStatus(id string) (*Domain.Loan, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	loan, err := u.loanRepository.GetLoanByID(objectID)
	if err != nil {
		return nil, err
	}

	// Log the retrieval of loan status
	log := &Domain.SystemLog{
		Timestamp: time.Now(),
		UserID:    loan.UserID,
		Action:    "Loan Status Retrieved",
		Details:   "User retrieved the status of a loan",
	}
	_ = u.logRepository.CreateLog(log)

	return loan, nil
}

func (u *loanUsecase) GetAllLoans(status string, order string) ([]Domain.Loan, error) {
	loans, err := u.loanRepository.GetAllLoans(status, order)
	if err != nil {
		return nil, err
	}

	// Log the retrieval of all loans (admin action)
	log := &Domain.SystemLog{
		Timestamp: time.Now(),
		Action:    "All Loans Retrieved",
		Details:   "Admin retrieved all loan applications",
	}
	_ = u.logRepository.CreateLog(log)

	return loans, nil
}

func (u *loanUsecase) ApproveRejectLoan(id string, status string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = u.loanRepository.UpdateLoanStatus(objectID, status)
	if err != nil {
		return err
	}

	// Log the loan status update (approve/reject)
	log := &Domain.SystemLog{
		Timestamp: time.Now(),
		Action:    "Loan Status Updated",
		Details:   "Admin updated the status of a loan to " + status,
	}
	_ = u.logRepository.CreateLog(log)

	return nil
}

func (u *loanUsecase) DeleteLoan(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	err = u.loanRepository.DeleteLoan(objectID)
	if err != nil {
		return err
	}

	// Log the loan deletion (admin action)
	log := &Domain.SystemLog{
		Timestamp: time.Now(),
		Action:    "Loan Deleted",
		Details:   "Admin deleted a loan application",
	}
	_ = u.logRepository.CreateLog(log)

	return nil
}
