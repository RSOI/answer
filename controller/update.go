package controller

import (
	"encoding/json"
	"fmt"

	"github.com/RSOI/answer/model"
	"github.com/RSOI/answer/utils"
	"github.com/RSOI/answer/view"
)

// MakeBestPATCH remove answer
func MakeBestPATCH(body []byte) (*model.Answer, error) {
	var err error

	var AnswerToUpdate model.Answer
	var UpdatedAnswer model.Answer
	isBest := true
	err = json.Unmarshal(body, &AnswerToUpdate)
	AnswerToUpdate.IsBest = &isBest
	if err != nil {
		utils.LOG(fmt.Sprintf("Broken body. Error: %s", err.Error()))
		return nil, err
	}

	err = view.ValidateMakeBestAnswer(AnswerToUpdate)
	if err != nil {
		utils.LOG(fmt.Sprintf("Validation error: %s", err.Error()))
		return nil, err
	}

	UpdatedAnswer, err = AnswerModel.UpdateAnswer(AnswerToUpdate)
	if err != nil {
		utils.LOG(fmt.Sprintf("Data error: %s", err.Error()))
		return nil, err
	}

	utils.LOG("Answer marked as best successfully")
	return &UpdatedAnswer, nil
}
