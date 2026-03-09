package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <username> <password>")
		return
	}

	username := os.Args[1]
	password := os.Args[2]

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	// Session Login
	loginURL := "https://springzabdesk.szabist-isb.edu.pk/VerifyLogin.asp?sid=974494673"
	loginData := url.Values{}
	loginData.Set("txtLoginName", username)
	loginData.Set("txtPassword", password)
	loginData.Set("txtCampus_Id", "1")

	resp, err := client.PostForm(loginURL, loginData)
	if err != nil {
		return
	}
	resp.Body.Close()

	// Submit form for CSC 2205 Operating Systems
	attendanceURL := "https://springzabdesk.szabist-isb.edu.pk/Student/QryCourseAttendance.asp?OptionName=View%20Attendance&sid=974494733"

	formData := url.Values{}
	formData.Set("txtFac", "Fakhar")
	formData.Set("txtSem", "16227")
	formData.Set("txtSec", "3")
	formData.Set("txtCou", "2726")

	resp, err = client.PostForm(attendanceURL, formData)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println("--- Course Attendance Detail Content ---")
	fmt.Println(string(body))
}
