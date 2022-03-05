package handler

import (
	"capstone-backend/DAO"
	"errors"
	"gorm.io/gorm"
	"log"
)

func AddCourse(course *DAO.Course) error {
	var courseFound DAO.Course
	err := DAO.DB().Where("course_code = ?", course.CourseCode).Find(&courseFound).Error
	log.Print(err)
	log.Print(courseFound)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// save to database
		error := DAO.DB().Create(&course).Error
		if error != nil {
			return &Error{500, "DB Error"}
		}
		return nil
	}

	if err != nil {
		return &Error{500, "DB Error"}
	}

	// The course code has existed in the database, update the row with the new request
	log.Printf("%d has already existed in databse, add course with req %+v\n", course.CourseCode, course)
	err = DAO.DB().Model(&courseFound).Updates(course).Error
	if err != nil {
		return &Error{500, "DB Error"}
	}
	return nil
}
