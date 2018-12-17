package view

import (
	"github.com/RSOI/answer/model"
	"github.com/RSOI/answer/ui"
)

// ValidateNewAnswer returns nil if all the required form values are passed
func ValidateNewAnswer(data model.Answer) error {
	if data.Content == nil ||
		*data.Content == "" ||
		data.AuthorID == 0 ||
		data.QuestionID == 0 {
		return ui.ErrFieldsRequired
	}
	return nil
}

// ValidateDeleteAnswer returns true if parameter to delete found
func ValidateDeleteAnswer(data model.Answer) (string, error) {
	if data.ID != 0 {
		return "id", nil
	}

	if data.QuestionID != 0 {
		return "question_id", nil
	}

	if data.AuthorID != 0 {
		return "author_id", nil
	}

	return "", ui.ErrFieldsRequired
}

// ValidateMakeBestAnswer returns nil if parameter to delete found
func ValidateMakeBestAnswer(data model.Answer) error {
	if data.ID != 0 {
		return nil
	}
	return ui.ErrFieldsRequired
}
