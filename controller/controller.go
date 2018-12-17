package controller

import (
	"github.com/RSOI/answer/model"
	"github.com/RSOI/answer/utils"
	"github.com/jackc/pgx"
)

var (
	// AnswerModel interface with methods
	AnswerModel model.AServiceInterface
)

// Init Init model with pgx connection
func Init(db *pgx.ConnPool) {
	utils.LOG("Setup model...")
	AnswerModel = &model.AService{
		Conn: db,
	}
}
