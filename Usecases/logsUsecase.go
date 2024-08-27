package Usecases

import (
	"LoanTracker/Domain"
	"LoanTracker/Repository"
)

type LogUsecase interface {
	GetLogs() ([]Domain.SystemLog, error)
}

type logUsecase struct {
	logRepository Repository.LogRepository
}

func NewLogUsecase(logRepo Repository.LogRepository) LogUsecase {
	return &logUsecase{
		logRepository: logRepo,
	}
}

func (u *logUsecase) GetLogs() ([]Domain.SystemLog, error) {
	return u.logRepository.GetLogs()
}
