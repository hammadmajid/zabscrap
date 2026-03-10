package scraper

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"zabscrap/internal/models"
)

// Scraper handles attendance data retrieval from ZabDesk
type Scraper struct {
	client *http.Client
}

// New creates a new Scraper instance
func New() *Scraper {
	jar, _ := cookiejar.New(nil)
	return &Scraper{
		client: &http.Client{Jar: jar},
	}
}

// ScrapeAttendance fetches and parses attendance data from ZabDesk
func (s *Scraper) ScrapeAttendance(username, password string) ([]models.CourseAttendance, error) {
	loginURL := "https://springzabdesk.szabist-isb.edu.pk/VerifyLogin.asp"
	resp, err := s.client.PostForm(loginURL, url.Values{
		"txtLoginName": {username},
		"txtPassword":  {password},
		"txtCampus_Id": {"1"},
	})
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	listURL := "https://springzabdesk.szabist-isb.edu.pk/Student/QryCourseAttendance.asp?OptionName=View%20Attendance"
	listResp, err := s.client.Get(listURL)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(listResp.Body)

	b, err := io.ReadAll(listResp.Body)
	if err != nil {
		return nil, err
	}

	reLinks := regexp.MustCompile(`chkSubmit\('([^']+)','([^']+)','([^']+)','([^']+)'\)`)
	matches := reLinks.FindAllStringSubmatch(string(b), -1)

	var results []models.CourseAttendance
	for _, m := range matches {
		dResp, err := s.client.PostForm(listURL, url.Values{
			"txtFac": {m[1]}, "txtSem": {m[2]}, "txtSec": {m[3]}, "txtCou": {m[4]},
		})
		if err != nil {
			return nil, err
		}

		db, err := io.ReadAll(dResp.Body)
		if err != nil {
			err := dResp.Body.Close()
			if err != nil {
				return nil, err
			}
			return nil, err
		}

		if err := dResp.Body.Close(); err != nil {
			return nil, err
		}

		dHTML := string(db)

		course := models.CourseAttendance{
			CourseName: s.parseTag(dHTML, "Course:"),
			Instructor: s.parseTag(dHTML, "Instructor:"),
		}

		reRow := regexp.MustCompile(`(?s)<tr>\s*<td[^>]*>(\d+)</td>\s*<td[^>]*>([\d/]+)</td>\s*<td[^>]*>\s*([a-zA-Z]+)\s*</td>\s*</tr>`)
		rows := reRow.FindAllStringSubmatch(dHTML, -1)
		for _, r := range rows {
			course.Records = append(course.Records, models.AttendanceRecord{
				Lecture: r[1], Date: r[2], Status: strings.TrimSpace(r[3]),
			})
		}
		results = append(results, course)
	}
	return results, nil
}

// parseTag extracts a value from an HTML table row by label
func (s *Scraper) parseTag(html, label string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?i)<th[^>]*>%s</th>\s*<td[^>]*>(.*?)</td>`, regexp.QuoteMeta(label)))
	m := re.FindStringSubmatch(html)
	if len(m) > 1 {
		return strings.TrimSpace(regexp.MustCompile("<[^>]*>").ReplaceAllString(m[1], ""))
	}
	return "N/A"
}
