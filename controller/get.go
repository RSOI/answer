package controller

import (
	"fmt"
	"strconv"

	"github.com/RSOI/answer/model"
	"github.com/RSOI/answer/utils"
)

// AnswerGET get answer by id
func AnswerGET(id string) (*model.Answer, error) {
	aID, _ := strconv.Atoi(id)

	data, err := AnswerModel.GetAnswerByID(aID)
	if err != nil {
		utils.LOG(fmt.Sprintf("Data error: %s", err.Error()))
		return nil, err
	}

	utils.LOG("Answer was found successfully")
	return &data, nil
}

// AnswersGET get answers by author or question
func AnswersGET(aid string, searchby string, limit int, offset int) ([]model.Answer, error) {
	var err error
	var data []model.Answer

	aidi, _ := strconv.Atoi(aid)
	switch searchby {
	case "author":
		data, err = AnswerModel.GetAnswersByAuthorID(aidi, limit, offset)
		break
	case "question":
		data, err = AnswerModel.GetAnswersByQuestionID(aidi, limit, offset)
		break
	}

	if err != nil {
		utils.LOG(fmt.Sprintf("Data error: %s", err.Error()))
		return nil, err
	}

	utils.LOG("Answers were found successfully")
	return data, nil
}
