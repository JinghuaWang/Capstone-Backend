package DAO

import (
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	CourseFullName string  `gorm:"size:255;not null"`
	CourseCode     string  `gorm:"size:50;uniqueIndex"`
	DepartmentCode string  `gorm:"size:50""`
	CourseNo       int     `gorm:"type:int"`
	CourseTile     string  `gorm:"size:255""`
	Credit         string  `gorm:"size:50"`
	Description    string  `gorm:"type:text"`
	Tags           Strings `gorm:"type:json"`
}
