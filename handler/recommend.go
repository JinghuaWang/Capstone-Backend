package handler

import "github.com/gin-gonic/gin"

var COURSES = RecommendCourses{
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
	Tech: Tech{
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
	DataResp(c, COURSES)
}
