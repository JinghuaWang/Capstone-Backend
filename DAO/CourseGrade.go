package DAO

import (
	"gorm.io/gorm"
)

type CourseGrade struct {
	gorm.Model
	CourseCode    string  `gorm:"size:255"`
	CourseId      uint    `gorm:"type:int"`
	Year          int     `gorm:"type:int"`
	Quarter       int     `gorm:"type:int"`
	ProfessorName string  `gorm:"size:255"`
	AverageGPA    float32 `gorm:"type:decimal(10,2);column:average_gap"`
	StudentCount  int     `gorm:"type:int"`
	A             int     `gorm:"type:int"`
	AMinus        int     `gorm:"type:int"`
	BPlus         int     `gorm:"type:int"`
	B             int     `gorm:"type:int"`
	BMinus        int     `gorm:"type:int"`
	CPlus         int     `gorm:"type:int"`
	C             int     `gorm:"type:int"`
	CMinus        int     `gorm:"type:int"`
	DPlus         int     `gorm:"type:int"`
	D             int     `gorm:"type:int"`
	DMinus        int     `gorm:"type:int"`
	Fail          int     `gorm:"type:int"`
	Withdraw      int     `gorm:"type:int"`
}
