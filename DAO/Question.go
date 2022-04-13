package DAO

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	CourseCode   string `gorm:"size:50;index"`
	QuestionText string `gorm:"type:text"`
	IPAddress    string `gorm:"size:50"`
	IsPreset     int    `gorm:"type:int"`
}
