package main

import (
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/RSOI/answer/controller"
	"github.com/RSOI/answer/model"
	"github.com/RSOI/answer/ui"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedAService struct {
	mock.Mock
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
func (s *MockedAService) GetAnswersByAuthorID(aAuthorID int) ([]model.Answer, error) {
	args := s.Mock.Called(aAuthorID)
	return args.Get(0).([]model.Answer), args.Error(1)
}
func (s *MockedAService) GetAnswersByQuestionID(aQuestionID int) ([]model.Answer, error) {
	args := s.Mock.Called(aQuestionID)
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
	HOST                        = "http://localhost"
	defaultAnswerContent        = "My Answer Content"
	defaultAnswerIsBest         = false
	updatedAnswerIsBest         = true
	defaultAnswerCreatedTime, _ = time.Parse("2006-01-02T15:04:05", time.Now().String())
	defaultAnswer               = model.Answer{
		AuthorID:   1,
		QuestionID: 1,
		Content:    &defaultAnswerContent,
	}
	createdAnswer = model.Answer{
		ID:         1,
		AuthorID:   1,
		QuestionID: 1,
		Content:    &defaultAnswerContent,
		IsBest:     &defaultAnswerIsBest,
		Created:    defaultAnswerCreatedTime,
	}
	updatedAnswer = model.Answer{
		ID:         1,
		AuthorID:   1,
		QuestionID: 1,
		Content:    &defaultAnswerContent,
		IsBest:     &updatedAnswerIsBest,
		Created:    defaultAnswerCreatedTime,
	}
)

func initServer() (*fasthttp.Client, *fasthttp.Request, *fasthttp.Response, *MockedAService) {
	listener := fasthttputil.NewInmemoryListener()
	server := &fasthttp.Server{
		Handler: initRoutes().Handler,
	}
	go server.Serve(listener)

	client := &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return listener.Dial()
		},
	}

	controller.AnswerModel = &MockedAService{}
	cMock := controller.AnswerModel.(*MockedAService)
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	return client, req, res, cMock
}

/*
********************************************************************
TESTS FOR ANSWER ***************************************************
********************************************************************
*/

func TestAnswerCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	a, _ := json.Marshal(&defaultAnswer)

	req.SetRequestURI(HOST + "/answer")
	req.Header.SetMethod("PUT")
	req.SetBody(a)

	cMock.On("AddAnswer", defaultAnswer).Return(createdAnswer, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 201, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 201, response.Status)
		assert.Equal(t, "", response.Error)
		responseData := response.Data.(map[string]interface{})
		assert.Equal(t, createdAnswer.ID, int(responseData["id"].(float64)))
		assert.Equal(t, *createdAnswer.Content, responseData["content"])
		assert.Equal(t, createdAnswer.AuthorID, int(responseData["author_id"].(float64)))
		assert.Equal(t, createdAnswer.QuestionID, int(responseData["question_id"].(float64)))
		assert.Equal(t, *createdAnswer.IsBest, responseData["is_best"])
		assert.Equal(t, createdAnswer.Created.Format("2006-01-02T15:04:05")+"Z", responseData["created"])
	}
}

func TestAnswerMissedOneField(t *testing.T) {
	client, req, res, _ := initServer()

	req.SetRequestURI(HOST + "/answer")
	req.Header.SetMethod("PUT")
	req.SetBodyString("{\"author_id\": 1, \"question_id\": 1}")

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		assert.Equal(t, 400, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 400, response.Status)
		assert.Equal(t, ui.ErrFieldsRequired.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestAnswerMissedManyField(t *testing.T) {
	client, req, res, _ := initServer()

	req.SetRequestURI(HOST + "/answer")
	req.Header.SetMethod("PUT")
	req.SetBodyString("{\"author_id\": 1}")

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		assert.Equal(t, 400, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 400, response.Status)
		assert.Equal(t, ui.ErrFieldsRequired.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestAnswerBrokenBody(t *testing.T) {
	client, req, res, _ := initServer()

	req.SetRequestURI(HOST + "/answer")
	req.Header.SetMethod("PUT")
	req.SetBodyString("{author_id: 1}")

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		assert.Equal(t, 500, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 500, response.Status)
		assert.NotEqual(t, "", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

/*
********************************************************************
TESTS FOR QUESTION ID **********************************************
********************************************************************
*/

func TestAnswerGetByIDCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	req.SetRequestURI(HOST + "/answer/id1")
	req.Header.SetMethod("GET")

	cMock.On("GetAnswerByID", 1).Return(createdAnswer, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)
		responseData := response.Data.(map[string]interface{})
		assert.Equal(t, createdAnswer.ID, int(responseData["id"].(float64)))
		assert.Equal(t, createdAnswer.QuestionID, int(responseData["question_id"].(float64)))
		assert.Equal(t, *createdAnswer.Content, responseData["content"])
		assert.Equal(t, createdAnswer.AuthorID, int(responseData["author_id"].(float64)))
		assert.Equal(t, *createdAnswer.IsBest, responseData["is_best"])
		assert.Equal(t, createdAnswer.Created.Format("2006-01-02T15:04:05")+"Z", responseData["created"])
	}
}

func TestAnswerGetByIDNotFound(t *testing.T) {
	client, req, res, cMock := initServer()

	req.SetRequestURI(HOST + "/answer/id0")
	req.Header.SetMethod("GET")

	cMock.On("GetAnswerByID", 0).Return(model.Answer{}, ui.ErrNoResult)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 404, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 404, response.Status)
		assert.Equal(t, ui.ErrNoResult.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

/*
********************************************************************
TESTS FOR ANSWER AUTHORID ****************************************
********************************************************************
*/

func TestAnswerGetByAuthorIDCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	req.SetRequestURI(HOST + "/answers/author1")
	req.Header.SetMethod("GET")

	createdAnswers := make([]model.Answer, 0)
	createdAnswers = append(createdAnswers, createdAnswer)
	createdAnswers = append(createdAnswers, createdAnswer)
	cMock.On("GetAnswersByAuthorID", 1).Return(createdAnswers, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)
		responseData := response.Data.([]interface{})
		assert.Equal(t, 2, len(responseData))
		for _, d := range responseData {
			data := d.(map[string]interface{})
			assert.Equal(t, createdAnswer.ID, int(data["id"].(float64)))
			assert.Equal(t, createdAnswer.QuestionID, int(data["question_id"].(float64)))
			assert.Equal(t, *createdAnswer.Content, data["content"])
			assert.Equal(t, createdAnswer.AuthorID, int(data["author_id"].(float64)))
			assert.Equal(t, *createdAnswer.IsBest, data["is_best"])
			assert.Equal(t, createdAnswer.Created.Format("2006-01-02T15:04:05")+"Z", data["created"])
		}
	}
}

func TestAnswerGetByAuthorIDNotFound(t *testing.T) {
	client, req, res, cMock := initServer()

	req.SetRequestURI(HOST + "/answers/author1")
	req.Header.SetMethod("GET")

	createdAnswers := make([]model.Answer, 0)
	cMock.On("GetAnswersByAuthorID", 1).Return(createdAnswers, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)
		responseData := response.Data.([]interface{})
		assert.Equal(t, 0, len(responseData))
		assert.Equal(t, make([]interface{}, 0), response.Data)
	}
}

/*
********************************************************************
TESTS FOR ANSWER QUESTIONID ****************************************
********************************************************************
*/

func TestAnswerGetByQuestionIDCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	req.SetRequestURI(HOST + "/answers/question1")
	req.Header.SetMethod("GET")

	createdAnswers := make([]model.Answer, 0)
	createdAnswers = append(createdAnswers, createdAnswer)
	createdAnswers = append(createdAnswers, createdAnswer)
	cMock.On("GetAnswersByQuestionID", 1).Return(createdAnswers, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)
		responseData := response.Data.([]interface{})
		assert.Equal(t, 2, len(responseData))
		for _, d := range responseData {
			data := d.(map[string]interface{})
			assert.Equal(t, createdAnswer.ID, int(data["id"].(float64)))
			assert.Equal(t, createdAnswer.QuestionID, int(data["question_id"].(float64)))
			assert.Equal(t, *createdAnswer.Content, data["content"])
			assert.Equal(t, createdAnswer.AuthorID, int(data["author_id"].(float64)))
			assert.Equal(t, *createdAnswer.IsBest, data["is_best"])
			assert.Equal(t, createdAnswer.Created.Format("2006-01-02T15:04:05")+"Z", data["created"])
		}
	}
}

func TestAnswerGetByQuestionIDNotFound(t *testing.T) {
	client, req, res, cMock := initServer()

	req.SetRequestURI(HOST + "/answers/question1")
	req.Header.SetMethod("GET")

	createdAnswers := make([]model.Answer, 0)
	cMock.On("GetAnswersByQuestionID", 1).Return(createdAnswers, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)
		responseData := response.Data.([]interface{})
		assert.Equal(t, 0, len(responseData))
		assert.Equal(t, make([]interface{}, 0), response.Data)
	}
}

/*
********************************************************************
TESTS FOR UPDATE QUESTION ******************************************
********************************************************************
*/

func TestUpdateCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	source, _ := json.Marshal(updatedAnswer)

	req.SetRequestURI(HOST + "/best")
	req.Header.SetMethod("PATCH")
	req.SetBody(source)

	cMock.On("UpdateAnswer", updatedAnswer).Return(updatedAnswer, nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)

		responseData := response.Data.(map[string]interface{})
		assert.Equal(t, *updatedAnswer.IsBest, responseData["is_best"])
	}
}

func TestUpdateNotFound(t *testing.T) {
	client, req, res, cMock := initServer()

	source, _ := json.Marshal(updatedAnswer)

	req.SetRequestURI(HOST + "/best")
	req.Header.SetMethod("PATCH")
	req.SetBody(source)

	cMock.On("UpdateAnswer", updatedAnswer).Return(model.Answer{}, ui.ErrNoDataToUpdate)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 404, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 404, response.Status)
		assert.Equal(t, ui.ErrNoDataToUpdate.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestUpdateMissedID(t *testing.T) {
	client, req, res, _ := initServer()

	req.SetRequestURI(HOST + "/best")
	req.Header.SetMethod("PATCH")
	req.SetBodyString("{\"is_best\": true}")

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		assert.Equal(t, 400, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 400, response.Status)
		assert.Equal(t, "missed required field(s)", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestUpdateBrokenBody(t *testing.T) {
	client, req, res, _ := initServer()

	req.SetRequestURI(HOST + "/best")
	req.Header.SetMethod("PATCH")
	req.SetBodyString("{id: 1}")

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		assert.Equal(t, 500, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 500, response.Status)
		assert.NotEqual(t, "", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

/*
********************************************************************
TESTS FOR REMOVE ANSWER ********************************************
********************************************************************
*/

func TestRemoveByIDCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	answerToRemove := model.Answer{
		ID: 1,
	}
	answerToRemoveJSON, _ := json.Marshal(&answerToRemove)

	req.SetRequestURI(HOST + "/delete")
	req.Header.SetMethod("DELETE")
	req.SetBody(answerToRemoveJSON)

	cMock.On("DeleteAnswerByID", answerToRemove).Return(nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestRemoveByAuthorIDCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	answerToRemove := model.Answer{
		AuthorID: 1,
	}
	answerToRemoveJSON, _ := json.Marshal(&answerToRemove)

	req.SetRequestURI(HOST + "/delete")
	req.Header.SetMethod("DELETE")
	req.SetBody(answerToRemoveJSON)

	cMock.On("DeleteAnswerByAuthorID", answerToRemove).Return(nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestRemoveByQuestionIDCorrectData(t *testing.T) {
	client, req, res, cMock := initServer()

	answerToRemove := model.Answer{
		QuestionID: 1,
	}
	answerToRemoveJSON, _ := json.Marshal(&answerToRemove)

	req.SetRequestURI(HOST + "/delete")
	req.Header.SetMethod("DELETE")
	req.SetBody(answerToRemoveJSON)

	cMock.On("DeleteAnswerByQuestionID", answerToRemove).Return(nil)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 200, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 200, response.Status)
		assert.Equal(t, "", response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestRemoveByIDNotFound(t *testing.T) {
	client, req, res, cMock := initServer()

	answerToRemove := model.Answer{
		ID: 1,
	}
	answerToRemoveJSON, _ := json.Marshal(&answerToRemove)

	req.SetRequestURI(HOST + "/delete")
	req.Header.SetMethod("DELETE")
	req.SetBody(answerToRemoveJSON)

	cMock.On("DeleteAnswerByID", answerToRemove).Return(ui.ErrNoDataToDelete)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 404, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 404, response.Status)
		assert.Equal(t, ui.ErrNoDataToDelete.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}

func TestRemoveMissedIDs(t *testing.T) {
	client, req, res, cMock := initServer()

	answerToRemove, _ := json.Marshal(&model.Answer{})

	req.SetRequestURI(HOST + "/delete")
	req.Header.SetMethod("DELETE")
	req.SetBody(answerToRemove)

	err := client.Do(req, res)
	if assert.Nil(t, err) {
		cMock.AssertExpectations(t)
		assert.Equal(t, 400, res.Header.StatusCode())

		var response ui.Response
		json.Unmarshal(res.Body(), &response)
		assert.Equal(t, 400, response.Status)
		assert.Equal(t, ui.ErrFieldsRequired.Error(), response.Error)
		assert.Equal(t, nil, response.Data)
	}
}
