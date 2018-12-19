package controller

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/RSOI/answer/model"
	"github.com/RSOI/answer/ui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedAService struct {
	mock.Mock
}

func getMock() *MockedAService {
	AnswerModel = &MockedAService{}
	return AnswerModel.(*MockedAService)
}

func (s *MockedAService) AddAnswer(a model.Answer) (model.Answer, error) {
	args := s.Mock.Called(a)
	return args.Get(0).(model.Answer), args.Error(1)
}
func (s *MockedAService) DeleteAnswerByID(a model.Answer) error {
	args := s.Mock.Called(a)
	return args.Error(0)
}
func (s *MockedAService) DeleteAnswerByAuthorID(a model.Answer) error {
	args := s.Mock.Called(a)
	return args.Error(0)
}
func (s *MockedAService) DeleteAnswerByQuestionID(a model.Answer) error {
	args := s.Mock.Called(a)
	return args.Error(0)
}
func (s *MockedAService) GetAnswerByID(qID int) (model.Answer, error) {
	args := s.Mock.Called(qID)
	return args.Get(0).(model.Answer), args.Error(1)
}
func (s *MockedAService) GetAnswersByAuthorID(aAuthorID int, limit int, offset int) ([]model.Answer, error) {
	args := s.Mock.Called(aAuthorID, limit, offset)
	return args.Get(0).([]model.Answer), args.Error(1)
}
func (s *MockedAService) GetAnswersByQuestionID(aQuestionID int, limit int, offset int) ([]model.Answer, error) {
	args := s.Mock.Called(aQuestionID, limit, offset)
	return args.Get(0).([]model.Answer), args.Error(1)
}
func (s *MockedAService) UpdateAnswer(a model.Answer) (model.Answer, error) {
	args := s.Mock.Called(a)
	return args.Get(0).(model.Answer), args.Error(1)
}
func (s *MockedAService) GetUsageStatistic(host string) (model.ServiceStatus, error) {
	args := s.Mock.Called(host)
	return args.Get(0).(model.ServiceStatus), args.Error(1)
}
func (s *MockedAService) LogStat(request []byte, responseStatus int, responseError string) {
	// nothing interesting here, just store data without affecting main thread
}

var (
	defaultAnswerContent        = "My Answer Content"
	defaultAnswerIsBest         = false
	updatedAnswerIsBest         = true
	defaultAnswerCreatedTime, _ = time.Parse("2006-01-02T15:04:05", time.Now().String())
	defaultAnswer               = model.Answer{
		AuthorID:       1,
		AuthorNickname: "Test",
		QuestionID:     1,
		Content:        &defaultAnswerContent,
	}
	createdAnswer = model.Answer{
		ID:             1,
		AuthorID:       1,
		AuthorNickname: "Test",
		QuestionID:     1,
		Content:        &defaultAnswerContent,
		IsBest:         &defaultAnswerIsBest,
		Created:        defaultAnswerCreatedTime,
	}
	updatedAnswer = model.Answer{
		ID:             1,
		AuthorID:       1,
		AuthorNickname: "Test",
		QuestionID:     1,
		Content:        &defaultAnswerContent,
		IsBest:         &updatedAnswerIsBest,
		Created:        defaultAnswerCreatedTime,
	}
	answerToRemoveID = model.Answer{
		ID: 1,
	}
	answerToRemoveAuthorID = model.Answer{
		AuthorID: 1,
	}
	answerToRemoveQuestionID = model.Answer{
		QuestionID: 1,
	}
)

/*
********************************************************************
TESTS FOR ASNWER ***************************************************
********************************************************************
*/

func TestAskAddCorrectData(t *testing.T) {
	body, _ := json.Marshal(&defaultAnswer)

	cMock := getMock()
	cMock.On("AddAnswer", defaultAnswer).Return(createdAnswer, nil)

	data, err := AnswerPUT(body)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, createdAnswer.ID, data.ID)
		assert.Equal(t, createdAnswer.QuestionID, data.QuestionID)
		assert.Equal(t, *createdAnswer.Content, *data.Content)
		assert.Equal(t, createdAnswer.AuthorID, data.AuthorID)
		assert.Equal(t, *createdAnswer.IsBest, *data.IsBest)
		assert.Equal(t, createdAnswer.Created, data.Created)
	}
}

func TestAnswerMissedField(t *testing.T) {
	body := []byte("{\"author_id\": 1}")

	data, err := AnswerPUT(body)
	assert.Equal(t, ui.ErrFieldsRequired, err)
	assert.Nil(t, data)
}

func TestAnswerBrokenBody(t *testing.T) {
	body := []byte("{author_id: 1}")

	data, err := AnswerPUT(body)
	assert.NotNil(t, err)
	assert.Nil(t, data)
}

/*
********************************************************************
TESTS FOR QUESTION ID **********************************************
********************************************************************
*/

func TestAnswerGetByIDCorrectData(t *testing.T) {
	cMock := getMock()
	cMock.On("GetAnswerByID", 1).Return(createdAnswer, nil)

	data, err := AnswerGET("1")
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, createdAnswer.ID, data.ID)
		assert.Equal(t, createdAnswer.QuestionID, data.QuestionID)
		assert.Equal(t, *createdAnswer.Content, *data.Content)
		assert.Equal(t, createdAnswer.AuthorID, data.AuthorID)
		assert.Equal(t, *createdAnswer.IsBest, *data.IsBest)
		assert.Equal(t, createdAnswer.Created, data.Created)
	}
}

func TestAnswerGetByIDNotFound(t *testing.T) {
	cMock := getMock()
	cMock.On("GetAnswerByID", 0).Return(model.Answer{}, ui.ErrNoResult)

	data, err := AnswerGET("0")
	if assert.NotNil(t, err) {
		cMock.AssertExpectations(t)

		assert.Nil(t, data)
		assert.Equal(t, ui.ErrNoResult, err)
	}
}

func TestAnswerGetByAuthorIDCorrectData(t *testing.T) {
	cMock := getMock()

	createdAnswers := make([]model.Answer, 0)
	createdAnswers = append(createdAnswers, createdAnswer)
	createdAnswers = append(createdAnswers, createdAnswer)
	cMock.On("GetAnswersByAuthorID", 1, -1, -1).Return(createdAnswers, nil)

	data, err := AnswersGET("1", "author", -1, -1)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 2, len(data))
		for _, d := range data {
			assert.Equal(t, createdAnswer.ID, d.ID)
			assert.Equal(t, createdAnswer.QuestionID, d.QuestionID)
			assert.Equal(t, *createdAnswer.Content, *d.Content)
			assert.Equal(t, createdAnswer.AuthorID, d.AuthorID)
			assert.Equal(t, *createdAnswer.IsBest, *d.IsBest)
			assert.Equal(t, createdAnswer.Created, d.Created)
		}
	}
}

func TestAnswerGetByAuthorIDNotFound(t *testing.T) {
	cMock := getMock()

	createdAnswers := make([]model.Answer, 0)
	cMock.On("GetAnswersByAuthorID", 1, -1, -1).Return(createdAnswers, nil)

	data, err := AnswersGET("1", "author", -1, -1)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, 0, len(data))
		assert.Equal(t, make([]model.Answer, 0), data)
	}
}

func TestAnswerGetByQuestionIDCorrectData(t *testing.T) {
	cMock := getMock()

	createdAnswers := make([]model.Answer, 0)
	createdAnswers = append(createdAnswers, createdAnswer)
	createdAnswers = append(createdAnswers, createdAnswer)
	cMock.On("GetAnswersByQuestionID", 1, -1, -1).Return(createdAnswers, nil)

	data, err := AnswersGET("1", "question", -1, -1)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 2, len(data))
		for _, d := range data {
			assert.Equal(t, createdAnswer.ID, d.ID)
			assert.Equal(t, createdAnswer.QuestionID, d.QuestionID)
			assert.Equal(t, *createdAnswer.Content, *d.Content)
			assert.Equal(t, createdAnswer.AuthorID, d.AuthorID)
			assert.Equal(t, *createdAnswer.IsBest, *d.IsBest)
			assert.Equal(t, createdAnswer.Created, d.Created)
		}
	}
}

func TestAnswerGetByQuestionIDNotFound(t *testing.T) {
	cMock := getMock()

	createdAnswers := make([]model.Answer, 0)
	cMock.On("GetAnswersByQuestionID", 1, -1, -1).Return(createdAnswers, nil)

	data, err := AnswersGET("1", "question", -1, -1)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, 0, len(data))
		assert.Equal(t, make([]model.Answer, 0), data)
	}
}

/*
********************************************************************
TESTS FOR UPDATE QUESTION ******************************************
********************************************************************
*/

func TestUpdateCorrectData(t *testing.T) {
	cMock := getMock()
	cMock.On("UpdateAnswer", updatedAnswer).Return(updatedAnswer, nil)

	body, _ := json.Marshal(updatedAnswer)
	response, err := MakeBestPATCH(body)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, *updatedAnswer.IsBest, *response.IsBest)
		assert.Equal(t, *updatedAnswer.Content, *response.Content)
	}
}

func TestUpdateNotFound(t *testing.T) {
	cMock := getMock()
	cMock.On("UpdateAnswer", updatedAnswer).Return(model.Answer{}, ui.ErrNoDataToUpdate)

	body, _ := json.Marshal(updatedAnswer)
	data, err := MakeBestPATCH(body)
	if assert.NotNil(t, err) {
		cMock.AssertExpectations(t)

		assert.Equal(t, ui.ErrNoDataToUpdate, err)
		assert.Nil(t, data)
	}
}

func TestUpdateMissedID(t *testing.T) {
	body := []byte("{\"has_best\": true, \"content\": \"My New Content\"}")

	response, err := MakeBestPATCH(body)
	assert.Equal(t, ui.ErrFieldsRequired, err)
	assert.Equal(t, (*model.Answer)(nil), response)
}

func TestUpdateBrokenBody(t *testing.T) {
	data, err := MakeBestPATCH([]byte("{id: 1}"))

	if assert.NotNil(t, err) {
		assert.Nil(t, data)
	}
}

/*
********************************************************************
TESTS FOR REMOVE QUESTION ******************************************
********************************************************************
*/

func TestRemoveByIDCorrectData(t *testing.T) {
	cMock := getMock()
	cMock.On("DeleteAnswerByID", answerToRemoveID).Return(nil)

	body := []byte("{\"id\": 1}")
	err := RemoveDELETE(body)
	assert.Nil(t, err)
}

func TestRemoveByAuthorIDCorrectData(t *testing.T) {
	cMock := getMock()
	cMock.On("DeleteAnswerByAuthorID", answerToRemoveAuthorID).Return(nil)

	body := []byte("{\"author_id\": 1}")
	err := RemoveDELETE(body)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
	}
}

func TestRemoveByQuestionIDCorrectData(t *testing.T) {
	cMock := getMock()
	cMock.On("DeleteAnswerByQuestionID", answerToRemoveQuestionID).Return(nil)

	body := []byte("{\"question_id\": 1}")
	err := RemoveDELETE(body)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
	}
}

func TestRemoveByIDNotFound(t *testing.T) {
	cMock := getMock()
	cMock.On("DeleteAnswerByID", answerToRemoveID).Return(ui.ErrNoDataToDelete)

	body := []byte("{\"id\": 1}")
	err := RemoveDELETE(body)
	if assert.Equal(t, ui.ErrNoDataToDelete, err) {
		cMock.AssertExpectations(t)
	}
}

func TestRemoveMissedIDs(t *testing.T) {
	body := []byte("{\"has_best\": true}")
	err := RemoveDELETE(body)
	assert.Equal(t, ui.ErrFieldsRequired, err)
}
