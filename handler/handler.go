package handler

import (
	"capstone-backend/DAO"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func IndexHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hi, welcome to the capstone backend index page!")
}

func GetCourseListHandler(c *gin.Context) {
	var resp = []CourseEntry{
		{
			CourseCode:     "CSE 142",
			CourseFullName: "CSE 142 Computer Programming I",
		},
		{
			CourseCode:     "CSE 143",
			CourseFullName: "CSE 143 Computer Programming II",
		},
		{
			CourseCode:     "MATH 126",
			CourseFullName: "MATH 126 Calculus With Analytic Geometry III",
		},
		{
			CourseCode:     "ENGL 131",
			CourseFullName: "ENGL 131 Composition: Exposition",
		},
	}

	DataResp(c, resp)
}

func GetCourseInfoHandler(c *gin.Context) {
	var description = "Basic programming-in-the-small abilities and concepts including procedural programming " +
		"(methods, parameters, return, values), basic control structures (sequence, if/else, for loop, while loop), " +
		"file processing, arrays, and an introduction to defining objects. Intended for students without prior " +
		"programming experience."
	var resp = CourseInfo{
		"CSE 142",
		"Computer Programming I",
		"4",
		[]string{"NW", "QSR"},
		description,
		[]string{"All", "Stuart Reges", "Brett Wortzman", "Miya Natsuhara"},
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

	token := strings.Split(req.CourseCode, " ")
	courseNo, err := strconv.Atoi(token[1])
	if len(token) < 2 || err != nil {
		InvalidResp(c, "Invalid Course Code (should in format of \"CSE 143\"")
		return
	}

	var courseInfo DAO.Course = DAO.Course{
		CourseFullName: req.CourseCode + " " + req.CourseTitle,
		CourseCode:     req.CourseCode,
		DepartmentCode: token[0],
		CourseNo:       courseNo,
		CourseTile:     req.CourseTitle,
		Credit:         req.Credit,
		Description:    req.Description,
		Tags:           req.Tags,
	}

	if err = AddCourse(&courseInfo); err != nil {
		ErrResp(c, err)
		return
	}
	SuccResp(c)
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
