package models

// AttendanceRecord represents a single attendance entry for a lecture
type AttendanceRecord struct {
	Lecture string `json:"lecture"`
	Date    string `json:"date"`
	Status  string `json:"status"`
}

// CourseAttendance represents attendance data for a course
type CourseAttendance struct {
	CourseName string             `json:"courseName"`
	Instructor string             `json:"instructor"`
	Records    []AttendanceRecord `json:"records"`
}
