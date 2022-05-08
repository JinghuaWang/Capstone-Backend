package handler

import "strconv"

// course list
type CourseEntry struct {
	CourseCode     string `json:"course_code"`
	CourseFullName string `json:"course_full_name"`
}

type CourseInfo struct {
	CourseCode  string   `json:"course_code"`
	CourseTitle string   `json:"course_title"`
	Credit      string   `json:"credit"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Professors  []string `json:"professors"`
}

type CourseSearch struct {
	QueryString string `json:"query_string"`
}

type CourseProfessorInfo struct {
	OverallRating     OneDecimal        `json:"overall_rating"`
	AverageGPA        float32           `json:"average_gpa"`
	Hours             OneDecimal        `json:"hours"`
	GradeDistribution GradeDistribution `json:"grade_distribution"`
	RatingBreakdown   RatingBreakdown   `json:"rating_breakdown"`
}

type GradeDistribution struct {
	A        int `json:"a"`
	AMinus   int `json:"a_minus"`
	BPlus    int `json:"b_plus"`
	B        int `json:"b"`
	BMinus   int `json:"b_minus"`
	CPlus    int `json:"c_plus"`
	C        int `json:"c"`
	CMinus   int `json:"c_minus"`
	DPlus    int `json:"d_plus"`
	D        int `json:"d"`
	DMinus   int `json:"d_minus"`
	Fail     int `json:"fail"`
	Withdraw int `json:"withdraw"`
}

type RatingBreakdown struct {
	InstructorContribution      OneDecimal `json:"instructor_contribution"`
	TeachingEffectiveness       OneDecimal `json:"teaching_effectiveness"`
	CourseOrganization          OneDecimal `json:"course_organization"`
	ClarityOfConceptExplanation OneDecimal `json:"clarity_of_concept_explanation"`
	AvailabilityOfExtraHelp     OneDecimal `json:"availability_of_extra_help"`
	UsefulnessOfCourseContent   OneDecimal `json:"usefulness_of_course_content"`
	GradingTechniques           OneDecimal `json:"grading_techniques"`
	ReasonableAssignedWork      OneDecimal `json:"reasonable_assigned_work"`
}

// Add course to course catalog
type CourseCreateInfoReq struct {
	CourseCode  string   `json:"course_code"`
	CourseTitle string   `json:"course_title"`
	Credit      string   `json:"credit"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
}

// Q&A
type QuestionInfo struct {
	QuestionID   int      `json:"question_id"`
	QuestionText string   `json:"question_text"`
	AskedAt      string   `json:"asked_at"`
	Answers      []Answer `json:"answers"`
}

type Answer struct {
	AnswerText string `json:"answer_text"`
	AnsweredAt string `json:"answered_at"`
}

type CreateQuestion struct {
	CourseCode   string `json:"course_code"`
	QuestionText string `json:"question_text"`
}

type CreateAnswer struct {
	CourseCode string `json:"course_code"`
	QuestionID int    `json:"question_id""`
	AnswerText string `json:"answer_text"`
}

type RecommendCourses struct {
	PopularCourses []*CourseBrief `json:"popular_courses"`
	GPABooster     []*CourseBrief `json:"gpa_booster"`
	Tech           Tech           `json:"tech"`
}

type Tech struct {
	DataAnalysisAndMachineLearning []*CourseBrief `json:"data_analysis_and_machine_learning"`
	Backend                        []*CourseBrief `json:"backend"`
	Frontend                       []*CourseBrief `json:"frontend"`
	Design                         []*CourseBrief `json:"design"`
}

type CourseBrief struct {
	CourseCode string   `json:"course_code"`
	Tags       []string `json:"tags"`
	Credit     string   `json:"credit"`
}

// helper struct
type OneDecimal float32

func (f OneDecimal) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(f), 'f', 1, 32)), nil
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Error struct {
	Code int
	Msg  string
}

func New(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

func (e *Error) Error() string {
	return e.Msg
}
