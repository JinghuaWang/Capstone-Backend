package handler

import (
	"capstone-backend/DAO"
	"fmt"
	"github.com/gin-gonic/gin"
)

// helper objects to configure recommendation courses
type recConfig struct {
	GPABooster     []string
	PopularCourses []string
	Tech           techConfig
}

type techConfig struct {
	DataAnalysisAndMachineLearning []string
	Backend                        []string
	Frontend                       []string
	Design                         []string
}

// recommendation course brief will be populated according the configuration
var recommendationCourse RecommendCourses

// configure recommandataion courses
var COURSES = recConfig{
	GPABooster: []string{
		"MUSIC 185",
		"ESS 101",
		"ARCH 150",
		"CLAS 101",
		"PSYCH 210",
		"DANCE 101",
		"DANCE 102",
		"INFO 101",
		"NUTR 400",
		"INFO 102",
		"MUSIC 116",
		"DANCE 100",
		"EDUC 210",
		"CLAS 430",
		"NUTR 200",
		"ECFS 200",
		"ESRM 100",
		"OCEAN 121",
		"DRAMA 101",
		"EDUC 215",
	},
	PopularCourses: []string{
		"CSE 142",
		"CSE 143",
		"BIOL 180",
		"MATH 124",
		"MATH 126",
		"ASTR 101",
		"COM 200",
		"ECON 200",
		"ECON 201",
		"PSYCH 101",
		"PSYCH 202",
		"CHEM 142",
		"PSYCH 210",
		"ENGL 131",
		"STAT 311",
		"INFO 200",
		"INFO 201",
		"CSE 373",
		"EDUC 251",
		"CLAS 430",
	},
	Tech: techConfig{
		DataAnalysisAndMachineLearning: []string{
			"INFO 201",
			"INFO 330",
			"CSE 414",
			"INFO 370",
			"INFO 371",
			"AMATH 482",
			"CSE 163",
			"CSE 416",
			"CSE 446",
		},
		Frontend: []string{
			"INFO 340",
			"CSE 154",
		},
		Backend: []string{
			"INFO 441",
			"CSE 331",
			"CSE 473",
			"CSE 154",
		},
		Design: []string{
			"DESIGN 166",
			"DESIGN 250",
			"HCDE 210",
			"INFO 200",
			"INFO 360",
		},
	},
}

func ListRecommendationHandler(c *gin.Context) {
	DataResp(c, recommendationCourse)
}

// initialize recommendation courses and cache into memory
func initRecommendationCourses() {
	// pre populate recommendationCourse
	recommendationCourse = RecommendCourses{
		PopularCourses: courseCodesToCourseBriefs(COURSES.PopularCourses),
		GPABooster:     courseCodesToCourseBriefs(COURSES.GPABooster),
		Tech: Tech{
			DataAnalysisAndMachineLearning: courseCodesToCourseBriefs(COURSES.Tech.DataAnalysisAndMachineLearning),
			Frontend:                       courseCodesToCourseBriefs(COURSES.Tech.Frontend),
			Backend:                        courseCodesToCourseBriefs(COURSES.Tech.Backend),
			Design:                         courseCodesToCourseBriefs(COURSES.Tech.Design),
		},
	}

	// get all the course code
	var courseCodes []string
	courseCodes = append(courseCodes, COURSES.PopularCourses...)
	courseCodes = append(courseCodes, COURSES.GPABooster...)
	courseCodes = append(courseCodes, COURSES.Tech.DataAnalysisAndMachineLearning...)
	courseCodes = append(courseCodes, COURSES.Tech.Frontend...)
	courseCodes = append(courseCodes, COURSES.Tech.Backend...)
	courseCodes = append(courseCodes, COURSES.Tech.Design...)

	// fetch course info
	var courses []DAO.Course
	err := DAO.DB().Where("course_code in ?", courseCodes).Find(&courses).Error
	if err != nil {
		panic(fmt.Sprintf("Fail to initialize course recommendation info %v", err))
	}

	// fill in other course info
	coursesToCourseBriefs(recommendationCourse.PopularCourses, courseMap)
	coursesToCourseBriefs(recommendationCourse.GPABooster, courseMap)
	coursesToCourseBriefs(recommendationCourse.Tech.DataAnalysisAndMachineLearning, courseMap)
	coursesToCourseBriefs(recommendationCourse.Tech.Frontend, courseMap)
	coursesToCourseBriefs(recommendationCourse.Tech.Backend, courseMap)
	coursesToCourseBriefs(recommendationCourse.Tech.Design, courseMap)
}

// populate course briefs by course codes
func courseCodesToCourseBriefs(courseCodes []string) []*CourseBrief {
	var briefs []*CourseBrief
	for _, courseCode := range courseCodes {
		briefs = append(briefs, &CourseBrief{
			CourseCode: courseCode,
		})
	}
	return briefs
}

// populate course info from course map
func coursesToCourseBriefs(briefs []*CourseBrief, courseMap map[string]DAO.Course) {
	for _, brief := range briefs {
		brief.Tags = courseMap[brief.CourseCode].Tags
		brief.Credit = courseMap[brief.CourseCode].Credit
	}
}
