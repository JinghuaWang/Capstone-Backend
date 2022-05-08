package dataimport

import (
	"capstone-backend/DAO"
	"fmt"
	"strconv"
	"strings"
)

// File path of the course evaluation median dataset
const CORUSE_EVAL_MEDIANS_FILE_PATH = "data/eval_medians2016-2021.csv"
const COURSE_PARAM_FILE_PATH = "data/evaluation_parameters2016-2021.csv"

var ACCEPTED_FORMS = map[string]bool{
	"A": true,
	"B": true,
	"C": true,
	"D": true,
	"E": true,
	"F": true,
	"G": true,
	"H": true,
}

func ImportCourseEvalData() {
	// read GPA data from local file
	courseEvals := readCSV(CORUSE_EVAL_MEDIANS_FILE_PATH)
	courseParam := readCSV(COURSE_PARAM_FILE_PATH)
	evalIDToCourseParameters := processCourseParams(courseParam)
	// translate CSV data to course evaluation data access objects
	courseEvalDAOs := processRawDataToDAO(courseEvals, evalIDToCourseParameters)

	// Insert into course evaluation table in size of 200 batches
	err := DAO.DB().CreateInBatches(courseEvalDAOs, 200).Error
	if err != nil {
		panic(fmt.Sprintf("Fail to insert course Evaluation data %v", err))
	}
}

func processCourseParams(courseParams [][]string) map[int]CourseParameter {
	evalIDToCourseParameters := make(map[int]CourseParameter)
	for i, param := range courseParams {
		// skip the header row
		if i == 0 {
			continue
		}
		evalID, _ := strconv.Atoi(param[3])
		numberOfResponse, _ := strconv.Atoi(param[8])
		CourseCode := param[17] + " " + param[18]
		section := param[19]
		evalIDToCourseParameters[evalID] = CourseParameter{
			EvalID:           evalID,
			CourseCode:       CourseCode,
			Section:          section,
			NumberOfResponse: numberOfResponse,
		}
	}
	return evalIDToCourseParameters
}

func processRawDataToDAO(courseEvals [][]string, evalIDToCourseParameters map[int]CourseParameter) []DAO.CourseEval {
	var courseEvalDAOs []DAO.CourseEval
	for i, eval := range courseEvals {
		// skip the header row
		if i == 0 {
			continue
		}
		form := strings.Trim(eval[4], " ")
		// skip not accepted forms
		if !ACCEPTED_FORMS[form] {
			continue
		}

		quarter := strings.Trim(eval[1], " ")
		year, _ := strconv.Atoi(eval[2])
		evalId, _ := strconv.Atoi(eval[3])
		overallRating, _ := strconv.ParseFloat(eval[5], 32)
		instructorsContribution, _ := strconv.ParseFloat(eval[8], 32)
		teachingEffectiveness, _ := strconv.ParseFloat(eval[9], 32)
		courseOrganization, _ := strconv.ParseFloat(eval[10], 32)
		clarityOfConceptExplanation, _ := strconv.ParseFloat(eval[12], 32)
		availabilityOfExtraHelp, _ := strconv.ParseFloat(eval[20], 32)
		usefulnessOfCourseContent, _ := strconv.ParseFloat(eval[20], 32)
		gradingTechniques, _ := strconv.ParseFloat(eval[25], 32)
		reasonableAssignedWork, _ := strconv.ParseFloat(eval[26], 32)
		hoursPerWeek, _ := strconv.ParseFloat(eval[33], 32)
		cei, _ := strconv.ParseFloat(eval[72], 32)

		evalDAO := DAO.CourseEval{
			CourseCode:                  evalIDToCourseParameters[evalId].CourseCode,
			Section:                     evalIDToCourseParameters[evalId].Section,
			Form:                        form,
			Year:                        year,
			Quarter:                     quarter,
			NumberOfResponse:            evalIDToCourseParameters[evalId].NumberOfResponse,
			ChallengeAndEngagementIndex: float32(cei),
			OverallRating:               float32(overallRating),
			HoursPerWeek:                float32(hoursPerWeek),
			InstructorsContribution:     float32(instructorsContribution),
			TeachingEffectiveness:       float32(teachingEffectiveness),
			CourseOrganization:          float32(courseOrganization),
			ClarityOfConceptExplanation: float32(clarityOfConceptExplanation),
			AvailabilityOfExtraHelp:     float32(availabilityOfExtraHelp),
			UsefulnessOfCourseContent:   float32(usefulnessOfCourseContent),
			GradingTechniques:           float32(gradingTechniques),
			ReasonableAssignedWork:      float32(reasonableAssignedWork),
		}
		courseEvalDAOs = append(courseEvalDAOs, evalDAO)
	}
	return courseEvalDAOs
}

// AggregateCourseEvaluation Aggregate course evaluation data by course code
// filters out quiz sections attached to lectures and evaluation with zero response.
func AggregateCourseEvaluation() {
	var courseEvalAggrs []DAO.CourseEvalAggr
	if err := DAO.
		DB().
		Table("course_evals").
		Select("course_code, sum(number_of_response) as number_of_response, " +
			"cast(sum(challenge_and_engagement_index * number_of_response) / sum(number_of_response) as decimal(10,2)) as challenge_and_engagement_index," +
			"cast(sum(overall_rating * number_of_response) / sum(number_of_response) as decimal(10,2)) as overall_rating," +
			"cast(sum(hours_per_week * number_of_response) / sum(number_of_response) as decimal(10,2)) as hours_per_week," +
			"cast(sum(instructors_contribution * number_of_response) / sum(number_of_response) as decimal(10,2)) as instructors_contribution," +
			"cast(sum(teaching_effectiveness * number_of_response) / sum(number_of_response) as decimal(10,2)) as teaching_effectiveness," +
			"cast(sum(course_organization * number_of_response) / sum(number_of_response) as decimal(10,2)) as course_organization," +
			"cast(sum(clarity_of_concept_explanation * number_of_response) / sum(number_of_response) as decimal(10,2)) as clarity_of_concept_explanation," +
			"cast(sum(availability_of_extra_help * number_of_response) / sum(number_of_response) as decimal(10,2)) as availability_of_extra_help," +
			"cast(sum(usefulness_of_course_content * number_of_response) / sum(number_of_response) as decimal(10,2)) as usefulness_of_course_content," +
			"cast(sum(grading_techniques * number_of_response) / sum(number_of_response) as decimal(10,2)) as grading_techniques," +
			"cast(sum(reasonable_assigned_work * number_of_response) / sum(number_of_response) as decimal(10,2)) as reasonable_assigned_work").
		Where("deleted_at is NULL AND number_of_response > 0 AND CHAR_LENGTH(section) = 1").
		Group("course_code").
		Find(&courseEvalAggrs).
		Error; err != nil {
		panic(fmt.Sprintf("Fail to aggreate course evaluation data, err: %v", err))
	}

	// insert the aggregated data into the course_evals_aggrs table
	if err := DAO.DB().CreateInBatches(&courseEvalAggrs, 100).Error; err != nil {
		panic(fmt.Sprintf("Fail to insert aggreated course evaluation data, err: %v", err))
	}
}

type CourseParameter struct {
	EvalID           int
	CourseCode       string
	Section          string
	NumberOfResponse int
}
