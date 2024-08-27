package Controller

import (
	"LoanTracker/Usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogController struct {
	LogUsecase Usecases.LogUsecase 
}

func NewLogController(usecase Usecases.LogUsecase) *LogController { 
	return &LogController{LogUsecase: usecase} 
}

func (lc *LogController) GetLogs(c *gin.Context) {
	logs, err := lc.LogUsecase.GetLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
