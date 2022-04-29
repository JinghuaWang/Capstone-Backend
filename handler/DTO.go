package handler

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
	OverallRating     float32           `json:"overall_rating"`
	AverageGPA        float32           `json:"average_gpa"`
	Hours             float32           `json:"hours"`
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
	TheCourseAsAWhole       float32 `json:"the_course_as_a_whole"`
	TheCourseContent        float32 `json:"the_course_content"`
	InstructorContribution  float32 `json:"instructor_contribution"`
	InstructorEffectiveness float32 `json:"instructor_effectiveness"`
	InstructorInterest      float32 `json:"instructor_interest"`
	QuizSectionContent      float32 `json:"quiz_section_content"`
	GradingTechniques       float32 `json:"grading_techniques"`
	AmountLearn             float32 `json:"amount_learn"`
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
	GPABooster     []string `json:"gpa_booster"`
	PopularCourses []string `json:"popular_courses"`
	Tech           Tech     `json:"tech"`
}

type Tech struct {
	DataAnalysisAndMachineLearning []string `json:"data_analysis_and_machine_learning"`
	Backend                        []string `json:"backend"`
	Frontend                       []string `json:"frontend"`
	Design                         []string `json:"design"`
}

// helper struct
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
