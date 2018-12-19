package model

import (
	"fmt"

	"github.com/jackc/pgx"

	"github.com/RSOI/answer/ui"
	"github.com/RSOI/answer/utils"
)

// AddAnswer add new answer
func (service *AService) AddAnswer(a Answer) (Answer, error) {
	var err error

	utils.LOG("Accessing database...")
	row := service.Conn.QueryRow(`
		INSERT INTO answer.answer
			(question_id, content, author_id, author_nickname) VALUES ($1, $2, $3, $4)
			RETURNING id, created, is_best
	`, a.QuestionID, a.Content, a.AuthorID, a.AuthorNickname)

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
		&a.AuthorNickname,
		&a.IsBest,
		&a.Created)

	return a, err
}

func (service *AService) getAnswers(query string, id int, limit int, offset int) ([]Answer, error) {
	var err error
	a := make([]Answer, 0)

	utils.LOG("Accessing database...")
	rows, err := service.Conn.Query(query, id, limit, offset)
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
			&ta.AuthorNickname,
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
func (service *AService) GetAnswersByAuthorID(aAuthorID int, limit int, offset int) ([]Answer, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 20 // default
	}
	utils.LOG(fmt.Sprintf(`SELECT * FROM answer.answer WHERE author_id = %d ORDER BY id ASC LIMIT %d OFFSET %d`, aAuthorID, limit, offset))
	return service.getAnswers(`SELECT * FROM answer.answer WHERE author_id = $1 ORDER BY id ASC LIMIT $2 OFFSET $3`, aAuthorID, limit, offset)
}

// GetAnswersByQuestionID get answer data by it's id
func (service *AService) GetAnswersByQuestionID(aQuestionID int, limit int, offset int) ([]Answer, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 20 // default
	}
	utils.LOG(fmt.Sprintf(`SELECT * FROM answer.answer WHERE question_id = %d ORDER BY id ASC LIMIT %d OFFSET %d`, aQuestionID, limit, offset))
	return service.getAnswers(`SELECT * FROM answer.answer WHERE question_id = $1 ORDER BY id ASC LIMIT $2 OFFSET $3`, aQuestionID, limit, offset)
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
	res, err := service.Conn.Exec(`UPDATE answer.answer SET is_best = true WHERE id = $1`, a.ID)
	if err == nil && res.RowsAffected() != 1 {
		err = ui.ErrNoDataToUpdate
	}
	*currentAnswerData.IsBest = true
	return currentAnswerData, err
}
