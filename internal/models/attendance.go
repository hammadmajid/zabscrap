package models

// AttendanceRecord represents a single attendance entry for a lecture
type AttendanceRecord struct {
	Lecture string
	Date    string
	Status  string
}

// CourseAttendance represents attendance data for a course
type CourseAttendance struct {
	CourseName string
	Instructor string
	Records    []AttendanceRecord
}
