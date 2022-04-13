package handler

import (
	"capstone-backend/DAO"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func AddQuestionHandler(c *gin.Context) {
	// get request param
	var req CreateQuestion
	err := c.ShouldBindJSON(&req)
	if err != nil {
		InvalidResp(c, "Invalid Parameters")
		return
	}

	exist, err := CourseExist(req.CourseCode)
	if err != nil {
		ErrResp(c, err)
		return
	}
	if !exist {
		InvalidResp(c, "Invalid Course Code")
		return
	}

	// insert a new question to the database
	question := DAO.Question{
		CourseCode:   req.CourseCode,
		QuestionText: req.QuestionText,
		IPAddress:    c.ClientIP(),
	}
	if err = DAO.DB().Create(&question).Error; err != nil {
		ErrResp(c, &Error{500, "DB Error"})
		return
	}

	SuccResp(c)
}

func AddAnswerHandler(c *gin.Context) {
	// get request param
	var req CreateAnswer
	err := c.ShouldBindJSON(&req)
	if err != nil {
		InvalidResp(c, "Invalid Parameters")
		return
	}

	exist, err := CourseExist(req.CourseCode)
	if err != nil {
		ErrResp(c, err)
		return
	}
	if !exist {
		InvalidResp(c, "Invalid Course Code")
		return
	}

	answer := DAO.Answer{
		QuestionId: uint(req.QuestionID),
		AnswerText: req.AnswerText,
		IPAddress:  c.ClientIP(),
	}
	if err = DAO.DB().Create(&answer).Error; err != nil {
		ErrResp(c, &Error{500, "DB Error"})
		return
	}

	SuccResp(c)
}

func ListQuestionsHandler(c *gin.Context) {
	var courseCode = c.Query("course_code")

	if courseCode == "" {
		InvalidResp(c, "Invalid Parameters")
		return
	}

	exist, err := CourseExist(courseCode)
	if err != nil {
		ErrResp(c, err)
		return
	}
	if !exist {
		InvalidResp(c, "Invalid Course Code")
		return
	}

	type QuestionAnswer struct {
		Id              int
		UpdatedAt       time.Time
		QuestionText    string
		IsPreset        bool
		AnswerText      string
		AnswerUpdatedAt time.Time
	}
	var results []QuestionAnswer

	if e := DAO.DB().Table("questions").
		Select("questions.id, questions.updated_at, questions.question_text, questions.is_preset, answers.answer_text, answers.updated_at as answer_updated_at").
		Joins("left join answers on answers.question_id = questions.id").
		Where("questions.course_code = ? AND questions.deleted_at is NULL AND answers.deleted_at is NULL ", courseCode).
		Find(&results); e.Error != nil {
		fmt.Printf("Error: %+v", e.Error)
		ErrResp(c, &Error{500, "DB Error"})
		return
	}

	// process the result
	var qMap = make(map[int]*QuestionInfo)
	for _, qa := range results {
		if _, ok := qMap[qa.Id]; !ok {
			qMap[qa.Id] = &QuestionInfo{
				QuestionID:   qa.Id,
				QuestionText: qa.QuestionText,
				AskedAt:      qa.UpdatedAt,
			}
		}

		question := qMap[qa.Id]
		question.Answers = append(question.Answers, Answer{
			AnswerText: qa.AnswerText,
			AnsweredAt: qa.AnswerUpdatedAt,
		})
	}

	var resp []QuestionInfo
	for _, question := range qMap {
		resp = append(resp, *question)
	}

	DataResp(c, resp)
}
