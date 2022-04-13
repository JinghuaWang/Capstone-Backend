package DAO

import "gorm.io/gorm"

type Answer struct {
	gorm.Model
	QuestionId uint   `gorm:"type:int"`
	AnswerText string `gorm:"type:text"`
	IPAddress  string `gorm:"size:50"`
}
