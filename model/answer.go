package model

import (
	"time"

	"github.com/jackc/pgx"

	"github.com/RSOI/answer/ui"
	"github.com/RSOI/answer/utils"
)

// Answer interface
type Answer struct {
	ID         int       `json:"id"`
	QuestionID int       `json:"question_id"`
	Content    *string   `json:"content"`
	AuthorID   int       `json:"author_id"`
	IsBest     *bool     `json:"is_best"`
	Created    time.Time `json:"created"`
}

// AddAnswer add new answer
func (service *AService) AddAnswer(a Answer) (Answer, error) {
	var err error

	utils.LOG("Accessing database...")
	row := service.Conn.QueryRow(`
		INSERT INTO answer.answer
			(question_id, content, author_id) VALUES ($1, $2, $3)
			RETURNING id, created, is_best
	`, a.QuestionID, a.Content, a.AuthorID)

	err = row.Scan(&a.ID, &a.Created, &a.IsBest)
	return a, err
}

// DeleteAnswerByID delete answer by id
func (service *AService) DeleteAnswerByID(a Answer) error {
	utils.LOG("Accessing database...")
	res, err := service.Conn.Exec(`DELETE FROM answer.answer WHERE id = $1`, a.ID)
	if err == nil && res.RowsAffected() != 1 {
		err = ui.ErrNoDataToDelete
	}
	return err
}

// DeleteAnswerByAuthorID delete answer by author id
func (service *AService) DeleteAnswerByAuthorID(a Answer) error {
	utils.LOG("Accessing database...")
	res, err := service.Conn.Exec(`DELETE FROM answer.answer WHERE author_id = $1`, a.AuthorID)
	if err == nil && res.RowsAffected() != 1 {
		err = nil
	}
	return err
}

// DeleteAnswerByQuestionID delete answer by author id
func (service *AService) DeleteAnswerByQuestionID(a Answer) error {
	utils.LOG("Accessing database...")
	res, err := service.Conn.Exec(`DELETE FROM answer.answer WHERE question_id = $1`, a.QuestionID)
	if err == nil && res.RowsAffected() != 1 {
		err = nil
	}
	return err
}

// GetAnswerByID get answer data by it's id
func (service *AService) GetAnswerByID(aID int) (Answer, error) {
	var err error
	var a Answer

	utils.LOG("Accessing database...")
	row := service.Conn.QueryRow(`SELECT * FROM answer.answer WHERE id = $1`, aID)

	err = row.Scan(
		&a.ID,
		&a.QuestionID,
		&a.Content,
		&a.AuthorID,
		&a.IsBest,
		&a.Created)

	return a, err
}

func (service *AService) getAnswers(query string, id int) ([]Answer, error) {
	var err error
	a := make([]Answer, 0)

	utils.LOG("Accessing database...")
	rows, err := service.Conn.Query(query, id)
	if err != nil {
		return a, err
	}

	for rows.Next() {
		var ta Answer
		err = rows.Scan(
			&ta.ID,
			&ta.QuestionID,
			&ta.Content,
			&ta.AuthorID,
			&ta.IsBest,
			&ta.Created)

		if err != nil {
			return a, err
		}

		a = append(a, ta)
	}

	return a, err
}

// GetAnswersByAuthorID get answer data by it's id
func (service *AService) GetAnswersByAuthorID(aAuthorID int) ([]Answer, error) {
	return service.getAnswers(`SELECT * FROM answer.answer WHERE author_id = $1 ORDER BY id ASC`, aAuthorID)
}

// GetAnswersByQuestionID get answer data by it's id
func (service *AService) GetAnswersByQuestionID(aQuestionID int) ([]Answer, error) {
	return service.getAnswers(`SELECT * FROM answer.answer WHERE question_id = $1 ORDER BY id ASC`, aQuestionID)
}

// UpdateAnswer Mark answer as best
func (service *AService) UpdateAnswer(a Answer) (Answer, error) {
	currentAnswerData, err := service.GetAnswerByID(a.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = ui.ErrNoDataToUpdate
		}
		return a, err
	}

	utils.LOG("Accessing database...")
	res, err := service.Conn.Exec(`UPDATE answer.answer SET has_best = true`)
	if err == nil && res.RowsAffected() != 1 {
		err = ui.ErrNoDataToUpdate
	}
	*currentAnswerData.IsBest = true
	return currentAnswerData, err
}
