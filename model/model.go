package model

import (
	"time"

	"github.com/jackc/pgx"
)

// Answer interface
type Answer struct {
	ID             int       `json:"id"`
	QuestionID     int       `json:"question_id"`
	Content        *string   `json:"content"`
	AuthorID       int       `json:"author_id"`
	AuthorNickname string    `json:"author_nickname"`
	IsBest         *bool     `json:"is_best"`
	Created        time.Time `json:"created"`
}

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
	GetAnswersByAuthorID(aAuthorID int, limit int, offset int) ([]Answer, error)
	GetAnswersByQuestionID(aQuestionID int, limit int, offset int) ([]Answer, error)
	UpdateAnswer(a Answer) (Answer, error)
	GetUsageStatistic(host string) (ServiceStatus, error)
	LogStat(request []byte, responseStatus int, responseError string)
}
