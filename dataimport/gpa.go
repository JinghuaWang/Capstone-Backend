package dataimport

import (
	"capstone-backend/DAO"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// File location of GPA dataset
const GPA_FILE_PATH = "gpa.csv"

func ImportGPAData() {
	// read GPA data from local file
	data := readCSV(GPA_FILE_PATH)
	// get the course_code to course_id mapping
	codeToId := getCourseCodeToIDMap()
	// translate CSV data to CourseGrade DAO
	grades := processGPAData(data, codeToId)

	// Insert the grades into course_grades table in size of 200 batches
	err := DAO.DB().CreateInBatches(grades, 200).Error
	if err != nil {
		panic(fmt.Sprintf("Fail to insert GPA data %v", err))
	}
}

// read CSV file from the given string and return a 2D array of string
func readCSV(file string) [][]string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return data
}

// Fetch all courses and create a course code to course id mapping
func getCourseCodeToIDMap() map[string]uint {
	codeToId := make(map[string]uint)
	var allCourses []DAO.Course
	err := DAO.DB().Find(&allCourses).Error
	if err != nil {
		panic(fmt.Sprintf("Fail to get course mapping %v", err))
	}

	for _, course := range allCourses {
		codeToId[course.CourseCode] = course.ID
	}
	return codeToId
}

// Convert 2D array of string into CourseGrade DAO
func processGPAData(data [][]string, codeToId map[string]uint) []DAO.CourseGrade {
	var grades []DAO.CourseGrade
	missingCode := make(map[string]int)
	for _, line := range data {
		// if this is a data entry
		if len(line) > 0 && strings.HasPrefix(line[0], "20") {

			// trim off the section code from the course number
			courseCode := strings.Trim(line[2][:(len(line[2])-2)], " ")
			year, _ := strconv.Atoi(line[1][0:4])
			quarter, _ := strconv.Atoi(line[1][4:5])
			count, _ := strconv.Atoi(line[4])
			a, _ := strconv.Atoi(line[5])
			am, _ := strconv.Atoi(line[6])
			bp, _ := strconv.Atoi(line[7])
			b, _ := strconv.Atoi(line[8])
			bm, _ := strconv.Atoi(line[9])
			cp, _ := strconv.Atoi(line[10])
			c, _ := strconv.Atoi(line[11])
			cm, _ := strconv.Atoi(line[12])
			dp, _ := strconv.Atoi(line[13])
			d, _ := strconv.Atoi(line[14])
			dm, _ := strconv.Atoi(line[15])
			f, _ := strconv.Atoi(line[16])
			w, _ := strconv.Atoi(line[17])
			averageGPA, _ := strconv.ParseFloat(line[18], 32)

			if _, ok := codeToId[courseCode]; !ok {
				missingCode[courseCode] = 1
			}

			grade := DAO.CourseGrade{
				CourseCode:   courseCode,
				CourseId:     codeToId[courseCode],
				Year:         year,
				Quarter:      quarter,
				StudentCount: count,
				AverageGPA:   float32(averageGPA),
				A:            a,
				AMinus:       am,
				BPlus:        bp,
				B:            b,
				BMinus:       bm,
				CPlus:        cp,
				C:            c,
				CMinus:       cm,
				DPlus:        dp,
				D:            d,
				DMinus:       dm,
				Fail:         f,
				Withdraw:     w,
			}
			grades = append(grades, grade)
		}
	}

	log.Printf("Missing course id: %d", len(missingCode))
	for k, _ := range missingCode {
		log.Printf("Could find course ID for course code %s", k)
	}
	return grades
}

// AggregateGPA aggregates raw grade data for each course
func AggregateGPA() {
	var gradeAgrrs []DAO.CourseGradeAggr
	if err := DAO.
		DB().
		Table("course_grades").
		Select("course_code, cast(sum(average_gap * student_count) / sum(student_count) as decimal(10,2)) as average_gpa, " +
			"sum(student_count) as student_count, sum(a) as a, sum(a_minus) as a_minus, " +
			"sum(b_plus) as b_plus, sum(b) as b, sum(b_minus) as b_minus, " +
			"sum(c_plus) as c_plus, sum(c) as c, sum(c_minus) as c_minus, " +
			"sum(d_plus) as d_plus, sum(d) as d, sum(d_minus) as d_minus, " +
			"sum(fail) as fail, sum(withdraw) as withdraw").
		Group("course_code").
		Find(&gradeAgrrs).
		Error; err != nil {
		panic(fmt.Sprintf("Fail to aggreate GPA data, err: %v", err))
	}

	// insert the aggregated data into the course_grade_aggrs table
	if err := DAO.DB().CreateInBatches(&gradeAgrrs, 100).Error; err != nil {
		panic(fmt.Sprintf("Fail to insert aggreated GPA data, err: %v", err))
	}
}
