package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type AttendanceRecord struct {
	Lecture string `json:"lecture"`
	Date    string `json:"date"`
	Status  string `json:"status"`
}

type CourseAttendance struct {
	CourseName string             `json:"course_name"`
	Instructor string             `json:"instructor"`
	Records    []AttendanceRecord `json:"records"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <username> <password>")
		return
	}

	username, password := os.Args[1], os.Args[2]
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// 1. Login
	loginURL := "https://springzabdesk.szabist-isb.edu.pk/VerifyLogin.asp"
	data := url.Values{
		"txtLoginName": {username},
		"txtPassword":  {password},
		"txtCampus_Id": {"1"},
	}
	resp, err := client.PostForm(loginURL, data)
	if err != nil {
		return
	}
	resp.Body.Close()

	// 2. Get Main Attendance List
	listURL := "https://springzabdesk.szabist-isb.edu.pk/Student/QryCourseAttendance.asp?OptionName=View%20Attendance"
	resp, err = client.Get(listURL)
	if err != nil {
		return
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	listHTML := string(bodyBytes)

	// Regex to find chkSubmit('fac','sem','sec','cou')
	reLinks := regexp.MustCompile(`chkSubmit\('([^']+)','([^']+)','([^']+)','([^']+)'\)`)
	matches := reLinks.FindAllStringSubmatch(listHTML, -1)

	var allAttendance []CourseAttendance

	// 3. Iterate through each course
	for _, match := range matches {
		fac, sem, sec, cou := match[1], match[2], match[3], match[4]

		formData := url.Values{
			"txtFac": {fac},
			"txtSem": {sem},
			"txtSec": {sec},
			"txtCou": {cou},
		}

		resp, err = client.PostForm(listURL, formData)
		if err != nil {
			continue
		}
		detailBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		detailHTML := string(detailBytes)

		// Parse Detail Data
		course := CourseAttendance{
			CourseName: extractTagContent(detailHTML, "Course:", 3),
			Instructor: extractTagContent(detailHTML, "Instructor:", 3),
		}

		// Extract Attendance Rows
		reRow := regexp.MustCompile(`(?s)<tr>\s*<td[^>]*>(\d+)</td>\s*<td[^>]*>([\d/]+)</td>\s*<td[^>]*>\s*([a-zA-Z]+)\s*</td>\s*</tr>`)
		rowMatches := reRow.FindAllStringSubmatch(detailHTML, -1)

		for _, rm := range rowMatches {
			course.Records = append(course.Records, AttendanceRecord{
				Lecture: rm[1],
				Date:    rm[2],
				Status:  strings.TrimSpace(rm[3]),
			})
		}

		allAttendance = append(allAttendance, course)
	}

	// 4. Output JSON
	finalJSON, _ := json.MarshalIndent(allAttendance, "", "  ")
	fmt.Println(string(finalJSON))
}

func extractTagContent(html, label string, colspan int) string {
	pattern := fmt.Sprintf(`(?i)<th[^>]*>%s</th>\s*<td[^>]*>(.*?)</td>`, regexp.QuoteMeta(label))
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		clean := regexp.MustCompile("<[^>]*>").ReplaceAllString(match[1], "")
		return strings.TrimSpace(clean)
	}
	return "Unknown"
}
