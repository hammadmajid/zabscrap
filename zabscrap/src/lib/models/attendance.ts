export interface AttendanceRecord {
	lecture: string;
	date: string;
	status: string;
}

export interface CourseAttendance {
	courseName: string;
	instructor: string;
	records: AttendanceRecord[];
}
