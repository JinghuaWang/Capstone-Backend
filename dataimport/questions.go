package dataimport

import (
	"capstone-backend/DAO"
	"capstone-backend/handler"
)

var PRESET_QUESTIONS = []string{
	"What’s the most important skill set you’ve learned from this course?",
	"Is the workload heavy? How many hours do you spend every week?",
	"What’s the key to succeeding in this class?",
	"Is your professor caring and easy to communicate with?",
	"How hard is it to get a 3.5+ in this class?",
}

func InitPresetQuestions() {
	courses := handler.AllCourse()

	var courseQuestions []DAO.Question
	for _, course := range courses {
		for _, question := range PRESET_QUESTIONS {
			courseQuestion := DAO.Question{
				CourseCode:   course.CourseCode,
				QuestionText: question,
				IsPreset:     1,
			}
			courseQuestions = append(courseQuestions, courseQuestion)
		}
	}

	if err := DAO.DB().CreateInBatches(&courseQuestions, 100).Error; err != nil {
		panic("Fail to insert preset questions to DB")
	}
}
