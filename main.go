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

	// Step 1: Login
	loginURL := "https://springzabdesk.szabist-isb.edu.pk/VerifyLogin.asp?sid=974494673"
	data := url.Values{}
	data.Set("txtLoginName", username)
	data.Set("txtPassword", password)
	data.Set("txtCampus_Id", "1")

	resp, err := client.PostForm(loginURL, data)
	if err != nil {
		fmt.Printf("Login Error: %v\n", err)
		return
	}
	resp.Body.Close()

	// Step 2: Access Attendance Page
	// Note: Use the SID from your previous successful redirect
	attendanceURL := "https://springzabdesk.szabist-isb.edu.pk/Student/QryCourseAttendance.asp?OptionName=View%20Attendance&sid=974494733"

	resp, err = client.Get(attendanceURL)
	if err != nil {
		fmt.Printf("Attendance Fetch Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read Error: %v\n", err)
		return
	}

	fmt.Println("--- Attendance Page Content ---")
	fmt.Println(string(body))
}
