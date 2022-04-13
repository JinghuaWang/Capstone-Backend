package handler

import (
	"capstone-backend/DAO"
	"errors"
	"gorm.io/gorm"
	"log"
)

func AddCourse(course *DAO.Course) error {
	var courseFound DAO.Course
	result := DAO.DB().Where("course_code = ?", course.CourseCode).First(&courseFound)
	if result.Error != nil {
		// If this course hasn't created, add this to the courses table
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			error := DAO.DB().Create(&course).Error
			if error != nil {
				return &Error{500, "DB Error"}
			}
			return nil
		}
		return &Error{500, "DB Error"}
	}

	// The course code has existed in the database, update the row with the new request
	log.Printf("%s has already existed in databse, add course with req %+v\n", course.CourseCode, course)
	err := DAO.DB().Model(&courseFound).Updates(course).Error
	if err != nil {
		return &Error{500, "DB Error"}
	}
	return nil
}

func CourseExist(courseCode string) (bool, error) {
	course := DAO.Course{
		CourseCode: courseCode,
	}
	if err := DAO.DB().Where(course).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, &Error{500, "DB Error"}
	}
	return true, nil
}
