package model

import (
	"github.com/jackc/pgx"
)

// AService connection holder
type AService struct {
	Conn *pgx.ConnPool
}

// AServiceInterface answer methods interface
type AServiceInterface interface {
	AddAnswer(a Answer) (Answer, error)
	DeleteAnswerByID(a Answer) error
	DeleteAnswerByAuthorID(a Answer) error
	DeleteAnswerByQuestionID(a Answer) error
	GetAnswerByID(aID int) (Answer, error)
	GetAnswersByAuthorID(aAuthorID int) ([]Answer, error)
	GetAnswersByQuestionID(aAuthorID int) ([]Answer, error)
	UpdateAnswer(a Answer) (Answer, error)
	GetUsageStatistic(host string) (ServiceStatus, error)
	LogStat(request []byte, responseStatus int, responseError string)
}
