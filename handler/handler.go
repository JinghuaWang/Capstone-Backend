package handler

import (
	"capstone-backend/DAO"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

const SEARCH_RESULT_MAX = 10

var courses []CourseEntry
var courseList []DAO.Course
var courseFullNameMap = make(map[string]CourseEntry)
var courseMap = make(map[string]DAO.Course)

func IndexHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hi, welcome to the capstone backend index page!")
}

func InitCachedData() {
	// fetch all courses from the DB and cache the list
	var allCourses []DAO.Course

	err := DAO.DB().Order("course_code").Find(&allCourses).Error
	if err != nil {
		panic(fmt.Sprintf("Fail to initialize course list cache %v", err))
	}

	courses = make([]CourseEntry, len(allCourses))
	courseList = make([]DAO.Course, len(allCourses))
	for i, c := range allCourses {
		courseEntry := CourseEntry{
			CourseCode:     c.CourseCode,
			CourseFullName: c.CourseFullName,
		}

		courses[i] = courseEntry
		courseList[i] = DAO.Course{
			CourseCode: strings.ToLower(c.CourseCode),
			// remove all white space for search
			CourseTitle:    strings.ReplaceAll(strings.ToLower(c.CourseCode), " ", ""),
			CourseFullName: strings.ToLower(c.CourseFullName),
		}
		courseFullNameMap[strings.ToLower(c.CourseFullName)] = courseEntry
		courseMap[c.CourseCode] = c
	}
}

func GetCourseListHandler(c *gin.Context) {
	DataResp(c, courses)
}

func SearchCourseHandler(c *gin.Context) {
	var keyword = c.Query("query")

	if keyword == "" {
		InvalidResp(c, "Invalid Parameters")
		return
	}

	results := searchKeyword(strings.ToLower(keyword))
	var courseResults = make([]CourseEntry, 0)
	for _, v := range results {
		courseResults = append(courseResults, v)
	}
	DataResp(c, courseResults)
}

func searchKeyword(word string) map[string]CourseEntry {
	resultMap := make(map[string]CourseEntry, 0)
	size := 0
	for _, c := range courseList {
		match := false
		// find matches
		if c.CourseCode == word || strings.HasPrefix(c.CourseCode, word) {
			match = true
		}
		if !match && (word == c.CourseTitle || strings.HasPrefix(c.CourseTitle, word)) {
			match = true
		}

		// add first X result to list
		if match {
			if size < SEARCH_RESULT_MAX {
				resultMap[c.CourseFullName] = courseFullNameMap[c.CourseFullName]
				size++
			} else {
				return resultMap
			}
		}
	}

	for _, c := range courseList {
		match := false
		// find matches
		if strings.Contains(c.CourseFullName, word) {
			match = true
		}

		// add first X result to list
		if match {
			if size < SEARCH_RESULT_MAX {
				resultMap[c.CourseFullName] = courseFullNameMap[c.CourseFullName]
				size++
			} else {
				return resultMap
			}
		}
	}
	return resultMap
}

func GetCourseInfoHandler(c *gin.Context) {
	var keyword = c.Query("course_code")

	if keyword == "" {
		InvalidResp(c, "Invalid Parameters")
		return
	}

	course := DAO.Course{
		CourseCode: keyword,
	}

	if err := DAO.DB().Where(course).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			InvalidResp(c, "Invalid Course Code")
			return
		}
		ErrResp(c, &Error{500, "DB Error"})
		return
	}

	professor := []string{"All", "Stuart Reges", "Brett Wortzman", "Miya Natsuhara"}
	var resp = CourseInfo{
		course.CourseCode,
		course.CourseTitle,
		course.Credit,
		course.Tags,
		course.Description,
		professor,
	}

	DataResp(c, resp)
}

func GetCourseProfessorInfoHandler(c *gin.Context) {
	var rating = RatingBreakdown{
		4.5,
		4.5,
		4.5,
		4.7,
		4.5,
		4.5,
		4.7,
		4.5,
	}

	var grade = GradeDistribution{
		20,
		48,
		37,
		17,
		12,
		9,
		7,
		6,
		4,
		3,
		3,
		5,
		7,
	}

	var resp = CourseProfessorInfo{
		4.6,
		3.2,
		17,
		grade,
		rating,
	}
	DataResp(c, resp)
}

func AddCourseInfoHandler(c *gin.Context) {
	var req CourseCreateInfoReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		InvalidResp(c, "Invalid Parameters")
		return
	}

	course, err := CourseInfoReqToCourse(req)
	if err != nil {
		ErrResp(c, err)
		return
	}

	if err = AddCourse(course); err != nil {
		ErrResp(c, err)
		return
	}
	SuccResp(c)
}

func CourseInfoReqToCourse(req CourseCreateInfoReq) (*DAO.Course, error) {
	//token := strings.Split(req.CourseCode, " ")
	//courseNo, err := strconv.Atoi(token[1])
	//if len(token) < 2 || err != nil {
	//	return nil, &Error{400, "Invalid Course Code (should in format of \"(B) CSE 143\""}
	//}

	// Split the course code into (campus) + department code and course number
	index := indexOfFirstDigit(req.CourseCode)
	first := req.CourseCode[:index]
	courseNo, err := strconv.Atoi(req.CourseCode[index:])
	if err != nil {
		return nil, &Error{400, "Course number isn't an integer"}
	}

	// check whether it contains the campus code
	campusCode := "S"
	departmentCode := first
	if len(first) >= 2 && (first[0] == 'B' || first[0] == 'T') && first[1] == ' ' {
		campusCode = string(first[0])
		departmentCode = first[1:]
	}

	var course DAO.Course = DAO.Course{
		CourseFullName: req.CourseCode + " " + req.CourseTitle,
		CourseCode:     req.CourseCode,
		CampusCode:     campusCode,
		DepartmentCode: departmentCode,
		CourseNo:       courseNo,
		CourseTitle:    req.CourseTitle,
		Credit:         req.Credit,
		Description:    req.Description,
		Tags:           req.Tags,
	}

	return &course, nil
}

func indexOfFirstDigit(s string) int {
	for i, c := range s {
		if _, err := strconv.Atoi(string(c)); err == nil {
			return i
		}
	}
	return 0
}

// Helper functions

func DataResp(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func ErrResp(c *gin.Context, e error) {
	Error, ok := e.(*Error)
	if !ok {
		c.JSON(400, Response{Code: 400, Message: "Error Occurred"})
		return
	}

	c.JSON(Error.Code, Response{Code: Error.Code, Message: Error.Msg})
}

func InvalidResp(c *gin.Context, msg string) {
	c.JSON(400, Response{Code: 400, Message: msg})
}

func SuccResp(c *gin.Context) {
	c.JSON(200, Response{Code: 200, Message: "Success"})
}
