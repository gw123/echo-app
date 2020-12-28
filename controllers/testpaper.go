package controllers

import (
	"encoding/json"
	"strconv"

	echoapp "github.com/gw123/echo-app"
	"github.com/labstack/echo"
)

type TestpaperController struct {
	echoapp.BaseController
	testpaperSvc echoapp.TestpaperService
}

func NewTestpaperController(testpaperSvc echoapp.TestpaperService) *TestpaperController {
	return &TestpaperController{
		testpaperSvc: testpaperSvc,
	}
}

type QuestionAndOptions struct {
	QuestionId []int    `json:"question_id"`
	Type       []string `json:"type"`
	Answer     []string `json:"answer"`
}

var ans []*QuestionAndOptions

func (testpaperCtrl *TestpaperController) SaveUserAnswer(ctx echo.Context) error {
	var answer = &QuestionAndOptions{}

	if err := ctx.Bind(answer); err != nil {
		return testpaperCtrl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	userAnswer := &echoapp.UserAnswer{}

	answerStr, err := json.Marshal(answer)
	if err != nil {
		return testpaperCtrl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	userAnswer.QAStr = string(answerStr)
	if err := testpaperCtrl.testpaperSvc.SaveUserTestAnswer(userAnswer); err != nil {
		return testpaperCtrl.Fail(ctx, echoapp.CodeNotFound, echoapp.ErrDb.Error(), err)
	}
	return testpaperCtrl.Success(ctx, answer)
}
func (testpaperCtrl *TestpaperController) SetTestpaper(ctx echo.Context) error {
	questionOp := &echoapp.Testpaper{}
	if err := ctx.Bind(questionOp); err != nil {
		return testpaperCtrl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	if err := testpaperCtrl.testpaperSvc.SaveTestpaper(questionOp); err != nil {
		return testpaperCtrl.Fail(ctx, echoapp.CodeNotFound, echoapp.ErrDb.Error(), err)
	}
	return testpaperCtrl.Success(ctx, questionOp)
}
func (testpaperCtrl *TestpaperController) GetTestpaperById(ctx echo.Context) error {
	id := ctx.QueryParam("id")
	idint, _ := strconv.ParseInt(id, 10, 64)
	test, err := testpaperCtrl.testpaperSvc.GetTestpaperById(idint)
	if err != nil {
		return testpaperCtrl.Fail(ctx, echoapp.CodeDBError, echoapp.ErrDb.Error(), err)
	}
	return testpaperCtrl.Success(ctx, test)
}
