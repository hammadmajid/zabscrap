package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

// Data Structures
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

// HTML Templates
const layoutHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="https://unpkg.com/terminal.css@0.7.4/dist/terminal.min.css" />
    <title>ZabDesk Scraper</title>
    <style>
        :root { --page-width: 800px; }
        .container { margin-top: 20px; }
        pre { background: #222; color: #00ff00; padding: 15px; border-radius: 5px; overflow-x: auto; }
    </style>
</head>
<body class="terminal">
    <div class="container">
        <header class="terminal-nav">
            <div class="terminal-logo">
                <div class="logo terminal-prompt"><a href="/" class="no-style">ZABDESK_SCRAPER_V1.0</a></div>
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
            <legend>Student Credentials</legend>
            <div class="form-group">
                <label for="username">Username:</label>
                <input id="username" name="username" type="text" required placeholder="e.g. 2312200">
            </div>
            <div class="form-group">
                <label for="password">Password:</label>
                <input id="password" name="password" type="password" required>
            </div>
            <div class="form-group">
                <button class="btn btn-default" type="submit" role="button">Execute Scraping</button>
            </div>
        </fieldset>
    </form>
</section>`

const resultHTML = `
<section>
    <h3>Extraction Results</h3>
    <pre><code>{{.JSON}}</code></pre>
    <a href="/" class="btn btn-primary">Back to Login</a>
</section>`

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/fetch", handleFetch)

	fmt.Println("Server initialized on :8080")
	http.ListenAndServe(":8080", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	render(w, formHTML, nil)
}

func handleFetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := r.FormValue("username")
	pass := r.FormValue("password")

	data, err := scrapeAttendance(user, pass)
	if err != nil {
		http.Error(w, fmt.Sprintf("Scraping failed: %v", err), http.StatusInternalServerError)
		return
	}

	prettyJSON, _ := json.MarshalIndent(data, "", "  ")
	render(w, resultHTML, struct{ JSON string }{string(prettyJSON)})
}

func render(w http.ResponseWriter, content string, data interface{}) {
	tmplString := strings.Replace(layoutHTML, "{{content}}", content, 1)
	tmpl, _ := template.New("page").Parse(tmplString)
	tmpl.Execute(w, data)
}

func scrapeAttendance(username, password string) ([]CourseAttendance, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// 1. Login
	loginURL := "https://springzabdesk.szabist-isb.edu.pk/VerifyLogin.asp?sid=974494673"
	vals := url.Values{
		"txtLoginName": {username},
		"txtPassword":  {password},
		"txtCampus_Id": {"1"},
	}
	resp, err := client.PostForm(loginURL, vals)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// 2. Access List
	listURL := "https://springzabdesk.szabist-isb.edu.pk/Student/QryCourseAttendance.asp?OptionName=View%20Attendance&sid=974494733"
	resp, err = client.Get(listURL)
	if err != nil {
		return nil, err
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	listHTML := string(bodyBytes)

	reLinks := regexp.MustCompile(`chkSubmit\('([^']+)','([^']+)','([^']+)','([^']+)'\)`)
	matches := reLinks.FindAllStringSubmatch(listHTML, -1)

	var results []CourseAttendance
	for _, m := range matches {
		formData := url.Values{
			"txtFac": {m[1]}, "txtSem": {m[2]}, "txtSec": {m[3]}, "txtCou": {m[4]},
		}
		dResp, err := client.PostForm(listURL, formData)
		if err != nil {
			continue
		}
		dBytes, _ := io.ReadAll(dResp.Body)
		dResp.Body.Close()
		dHTML := string(dBytes)

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
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return strings.TrimSpace(regexp.MustCompile("<[^>]*>").ReplaceAllString(match[1], ""))
	}
	return "N/A"
}
