package DAO

import (
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	CourseFullName string  `gorm:"size:255;not null"`   // B CSE 142 Computer Programming I
	CourseCode     string  `gorm:"size:50;uniqueIndex"` // B CSE 142
	CampusCode     string  `gorm:"size:10""`            // [S, B, T]
	DepartmentCode string  `gorm:"size:50""`            // CSE
	CourseNo       int     `gorm:"type:int"`            // 142
	CourseTitle    string  `gorm:"size:255""`
	Credit         string  `gorm:"size:50"`
	Description    string  `gorm:"type:text"`
	Tags           Strings `gorm:"type:json"`
}
