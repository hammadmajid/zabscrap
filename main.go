package main

import (
	"fmt"
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

	// Initialize cookie jar to handle session state
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginURL := "https://springzabdesk.szabist-isb.edu.pk/VerifyLogin.asp?sid=974494673"

	// Define form data based on HTML input names
	data := url.Values{}
	data.Set("txtLoginName", username)
	data.Set("txtPassword", password)
	data.Set("txtCampus_Id", "1") // Defaulted to Islamabad per HTML source

	// Execute POST request
	resp, err := client.PostForm(loginURL, data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check response location or body to verify success
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Final URL: %s\n", resp.Request.URL.String())
}
