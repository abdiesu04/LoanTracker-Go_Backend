package Controller

import (
	"LoanTracker/Domain"
	"LoanTracker/Usecases"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



type LoanController struct {
	LoanUsecase Usecases.LoanUsecase
}

func NewLoanController(usecase Usecases.LoanUsecase) *LoanController {
	return &LoanController{LoanUsecase: usecase}
}


func (c *LoanController) ApplyLoanHandler(ctx *gin.Context) {
	var loanRequest Domain.ApplyLoanRequest
	if err := ctx.ShouldBindJSON(&loanRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	response, err := c.LoanUsecase.ApplyLoan(&loanRequest, userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}


func (c *LoanController) GetLoanStatusHandler(ctx *gin.Context) {
	loanID := ctx.Param("id")
	// objectID, err := primitive.ObjectIDFromHex(loanID)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
	// 	return
	// }

	loan, err := c.LoanUsecase.GetLoanStatus(loanID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	ctx.JSON(http.StatusOK, loan)
}


func (c *LoanController) GetAllLoansHandler(ctx *gin.Context) {
	status := ctx.DefaultQuery("status", "all")
	order := ctx.DefaultQuery("order", "asc")

	loans, err := c.LoanUsecase.GetAllLoans(status, order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, loans)
}


func (c *LoanController) UpdateLoanStatusHandler(ctx *gin.Context) {
	loanID := ctx.Param("id")
	status := ctx.Param("status")
	if strings.ToLower(status) != "approved" && strings.ToLower(status) != "rejected" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	_, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	err = c.LoanUsecase.ApproveRejectLoan(loanID, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Loan status updated successfully"})
}


func (c *LoanController) DeleteLoanHandler(ctx *gin.Context) {
	loanID := ctx.Param("id")
	_, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	err = c.LoanUsecase.DeleteLoan(loanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Loan deleted successfully"})
}
