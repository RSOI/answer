package controller

import (
	"encoding/json"
	"fmt"

	"github.com/RSOI/answer/model"
	"github.com/RSOI/answer/utils"
	"github.com/RSOI/answer/view"
)

// AnswerPUT new answer
func AnswerPUT(body []byte) (*model.Answer, error) {
	var err error

	var NewAnswer model.Answer
	err = json.Unmarshal(body, &NewAnswer)
	if err != nil {
		utils.LOG(fmt.Sprintf("Broken body. Error: %s", err.Error()))
		return nil, err
	}

	err = view.ValidateNewAnswer(NewAnswer)
	if err != nil {
		utils.LOG(fmt.Sprintf("Validation error: %s", err.Error()))
		return nil, err
	}

	NewAnswer, err = AnswerModel.AddAnswer(NewAnswer)
	if err != nil {
		utils.LOG(fmt.Sprintf("Data error: %s", err.Error()))
		return nil, err
	}

	utils.LOG("New answer added successfully")
	return &NewAnswer, nil
}
