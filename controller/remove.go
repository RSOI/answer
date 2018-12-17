package controller

import (
	"encoding/json"
	"fmt"

	"github.com/RSOI/answer/model"
	"github.com/RSOI/answer/utils"
	"github.com/RSOI/answer/view"
)

// RemoveDELETE remove answer
func RemoveDELETE(body []byte) error {
	var err error

	var AnswerToRemove model.Answer
	err = json.Unmarshal(body, &AnswerToRemove)
	if err != nil {
		utils.LOG(fmt.Sprintf("Broken body. Error: %s", err.Error()))
		return err
	}

	f, err := view.ValidateDeleteAnswer(AnswerToRemove)
	if err != nil {
		utils.LOG(fmt.Sprintf("Validation error: %s", err.Error()))
		return err
	}

	utils.LOG(fmt.Sprintf("Removing answer by: %s...", f))

	switch f {
	case "id":
		err = AnswerModel.DeleteAnswerByID(AnswerToRemove)
	case "question_id":
		err = AnswerModel.DeleteAnswerByQuestionID(AnswerToRemove)
	case "author_id":
		err = AnswerModel.DeleteAnswerByAuthorID(AnswerToRemove)
	}

	if err != nil {
		utils.LOG(fmt.Sprintf("Data error: %s", err.Error()))
	} else {
		utils.LOG("Answer removed successfully")
	}

	return err
}
