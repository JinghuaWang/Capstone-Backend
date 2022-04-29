package handler

import (
	"capstone-backend/DAO"
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

	if _, ok := courseMap[req.CourseCode]; !ok {
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

	if _, ok := courseMap[req.CourseCode]; !ok {
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
	// input validation
	var courseCode = c.Query("course_code")

	if courseCode == "" {
		InvalidResp(c, "Invalid Parameters")
		return
	}

	// check the course code exist
	if _, ok := courseMap[courseCode]; !ok {
		InvalidResp(c, "Invalid Course Code")
		return
	}

	// Get the questions and answers by joining tables
	var results []QuestionAnswer
	if err := DAO.DB().Table("questions").
		Select("questions.id, questions.updated_at, questions.question_text, questions.is_preset, answers.answer_text, answers.updated_at as answer_updated_at").
		Joins("left join answers on answers.question_id = questions.id").
		Where("questions.course_code = ? AND questions.deleted_at is NULL AND answers.deleted_at is NULL ", courseCode).
		Order("questions.updated_at DESC, questions.id ASC, answers.updated_at DESC").
		Find(&results).
		Error; err != nil {
		ErrResp(c, &Error{500, "DB Error"})
		return
	}

	DataResp(c, qaListToNestedList(results))
}

// pre condition: questions are grouped by question id
func qaListToNestedList(list []QuestionAnswer) []*QuestionInfo {
	// group results by question id
	var resp []*QuestionInfo
	var prevID = 0
	var question *QuestionInfo
	for _, qa := range list {
		if prevID != qa.Id {
			prevID = qa.Id
			question = &QuestionInfo{
				QuestionID:   qa.Id,
				QuestionText: qa.QuestionText,
				AskedAt:      qa.UpdatedAt.Format("2006-01-02"),
			}
			resp = append(resp, question)
		}

		question.Answers = append(question.Answers, Answer{
			AnswerText: qa.AnswerText,
			AnsweredAt: qa.AnswerUpdatedAt.Format("2006-01-02"),
		})
	}

	return resp
}

type QuestionAnswer struct {
	Id              int
	UpdatedAt       time.Time
	QuestionText    string
	IsPreset        bool
	AnswerText      string
	AnswerUpdatedAt time.Time
}
