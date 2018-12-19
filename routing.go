package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/RSOI/answer/controller"
	"github.com/RSOI/answer/ui"
	"github.com/RSOI/answer/utils"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func sendResponse(ctx *fasthttp.RequestCtx, r ui.Response, nolog ...bool) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(r.Status)
	utils.LOG(fmt.Sprintf("Sending response. Status: %d", r.Status))

	doLog := true
	if len(nolog) > 0 {
		doLog = !nolog[0]
	}

	if doLog {
		controller.LogStat(ctx.Path(), r.Status, r.Error)
	}

	content, _ := json.Marshal(r)
	ctx.Write(content)
}

func indexGET(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Request: Get service stats (%s)", ctx.Path()))
	var err error
	var r ui.Response

	r.Data, err = controller.IndexGET(ctx.Host())
	r.Status, r.Error = ui.ErrToResponse(err)

	nolog := true
	sendResponse(ctx, r, nolog)
}

func answerPUT(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Answer question (%s)", ctx.Path()))
	var err error
	var r ui.Response

	r.Data, err = controller.AnswerPUT(ctx.PostBody())
	r.Status, r.Error = ui.ErrToResponse(err)
	if r.Status == 200 {
		r.Status = 201 // REST :)
	}
	sendResponse(ctx, r)
}

func answerGET(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Get one answer (%s)", ctx.Path()))
	var err error
	var r ui.Response

	id := ctx.UserValue("id").(string)
	r.Data, err = controller.AnswerGET(id)
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func answersAuthorGET(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Request: Get answers by author id (%s)", ctx.Path()))
	var err error
	var r ui.Response

	aid := ctx.UserValue("authorid").(string)
	limit := ctx.QueryArgs().Peek("limit")
	offset := ctx.QueryArgs().Peek("limit")
	l := 0
	o := 0
	if len(limit) > 0 && len(offset) > 0 {
		l, err = strconv.Atoi(string(limit))
		o, err = strconv.Atoi(string(offset))
	}
	r.Data, err = controller.AnswersGET(aid, "author", l, o)
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func answersQuestionGET(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Request: Get answers by question id (%s)", ctx.Path()))
	var err error
	var r ui.Response

	qid := ctx.UserValue("questionid").(string)
	page := ctx.QueryArgs().Peek("page")
	countOnPage := ctx.QueryArgs().Peek("conp")
	p := 0
	c := 0
	if len(page) > 0 && len(countOnPage) > 0 {
		p, err = strconv.Atoi(string(page))
		c, err = strconv.Atoi(string(countOnPage))
	}
	r.Data, err = controller.AnswersGET(qid, "question", c, (p-1)*c)
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func makeBestPATCH(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Request: Mark answer as best (%s)", ctx.Path()))
	var err error
	var r ui.Response

	r.Data, err = controller.MakeBestPATCH(ctx.PostBody())
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func removeDELETE(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Request: Delete answer (%s)", ctx.Path()))
	var err error
	var r ui.Response

	err = controller.RemoveDELETE(ctx.PostBody())
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func initRoutes() *fasthttprouter.Router {
	utils.LOG("Setup router...")
	router := fasthttprouter.New()
	router.GET("/", indexGET)
	router.PUT("/answer", answerPUT)
	router.GET("/answer/id:id", answerGET)
	router.GET("/answers/author:authorid", answersAuthorGET)
	router.GET("/answers/question:questionid", answersQuestionGET)
	router.PATCH("/best", makeBestPATCH)
	router.DELETE("/delete", removeDELETE)

	return router
}
