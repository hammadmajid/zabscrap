package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

type AttendanceRecord struct {
	Lecture string
	Date    string
	Status  string
}

type CourseAttendance struct {
	CourseName string
	Instructor string
	Records    []AttendanceRecord
}

const layoutHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="https://unpkg.com/terminal.css@0.7.4/dist/terminal.min.css" />
    <title>Attendance Extraction</title>
    <style>
        :root { --page-width: 900px; }
        .status-absent { color: #ff5555; font-weight: bold; }
        .status-present { color: #50fa7b; }
        .course-card { margin-bottom: 40px; border: 1px solid #444; padding: 20px; }
    </style>
</head>
<body class="terminal">
    <div class="container">
        <header class="terminal-nav">
            <div class="terminal-logo">
                <div class="logo terminal-prompt">DATA_EXTRACT_ATTENDANCE</div>
            </div>
        </header>
        {{content}}
    </div>
</body>
</html>`

const formHTML = `
<section>
    <form action="/fetch" method="POST">
        <fieldset>
            <legend>ZabDesk Access</legend>
            <div class="form-group">
                <label for="username">Username:</label>
                <input id="username" name="username" type="text" required>
            </div>
            <div class="form-group">
                <label for="password">Password:</label>
                <input id="password" name="password" type="password" required>
            </div>
            <div class="form-group">
                <button class="btn btn-default" type="submit">Parse Attendance</button>
            </div>
        </fieldset>
    </form>
</section>`

const resultHTML = `
<section>
    {{range .}}
    <div class="course-card">
        <header>
            <h2 style="margin-bottom: 0;">{{.CourseName}}</h2>
            <p style="color: #888;">Instructor: {{.Instructor}}</p>
        </header>
        <table class="terminal-table">
            <thead>
                <tr>
                    <th>Lec #</th>
                    <th>Date</th>
                    <th>Status</th>
                </tr>
            </thead>
            <tbody>
                {{range .Records}}
                <tr>
                    <td>{{.Lecture}}</td>
                    <td>{{.Date}}</td>
                    <td class="{{if eq .Status "Absent"}}status-absent{{else}}status-present{{end}}">
                        {{.Status}}
                    </td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>
    {{end}}
    <hr>
    <a href="/" class="btn btn-primary">Return</a>
</section>`

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, formHTML, nil)
	})
	http.HandleFunc("/fetch", handleFetch)
	http.ListenAndServe(":8080", nil)
}

func handleFetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	user, pass := r.FormValue("username"), r.FormValue("password")
	data, _ := scrapeAttendance(user, pass)
	render(w, resultHTML, data)
}

func render(w http.ResponseWriter, contentTmpl string, data interface{}) {
	tString := strings.Replace(layoutHTML, "{{content}}", contentTmpl, 1)
	tmpl, _ := template.New("v").Parse(tString)
	tmpl.Execute(w, data)
}

func scrapeAttendance(username, password string) ([]CourseAttendance, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	loginURL := "https://springzabdesk.szabist-isb.edu.pk/VerifyLogin.asp?sid=974494673"
	client.PostForm(loginURL, url.Values{
		"txtLoginName": {username},
		"txtPassword":  {password},
		"txtCampus_Id": {"1"},
	})

	listURL := "https://springzabdesk.szabist-isb.edu.pk/Student/QryCourseAttendance.asp?OptionName=View%20Attendance&sid=974494733"
	resp, _ := client.Get(listURL)
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	reLinks := regexp.MustCompile(`chkSubmit\('([^']+)','([^']+)','([^']+)','([^']+)'\)`)
	matches := reLinks.FindAllStringSubmatch(string(b), -1)

	var results []CourseAttendance
	for _, m := range matches {
		dResp, _ := client.PostForm(listURL, url.Values{
			"txtFac": {m[1]}, "txtSem": {m[2]}, "txtSec": {m[3]}, "txtCou": {m[4]},
		})
		db, _ := io.ReadAll(dResp.Body)
		dResp.Body.Close()
		dHTML := string(db)

		course := CourseAttendance{
			CourseName: parseTag(dHTML, "Course:"),
			Instructor: parseTag(dHTML, "Instructor:"),
		}

		reRow := regexp.MustCompile(`(?s)<tr>\s*<td[^>]*>(\d+)</td>\s*<td[^>]*>([\d/]+)</td>\s*<td[^>]*>\s*([a-zA-Z]+)\s*</td>\s*</tr>`)
		rows := reRow.FindAllStringSubmatch(dHTML, -1)
		for _, r := range rows {
			course.Records = append(course.Records, AttendanceRecord{
				Lecture: r[1], Date: r[2], Status: strings.TrimSpace(r[3]),
			})
		}
		results = append(results, course)
	}
	return results, nil
}

func parseTag(html, label string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?i)<th[^>]*>%s</th>\s*<td[^>]*>(.*?)</td>`, regexp.QuoteMeta(label)))
	m := re.FindStringSubmatch(html)
	if len(m) > 1 {
		return strings.TrimSpace(regexp.MustCompile("<[^>]*>").ReplaceAllString(m[1], ""))
	}
	return "N/A"
}
