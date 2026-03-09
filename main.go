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
		// http.Client follows redirects automatically
	}

	loginURL := "https://springzabdesk.szabist-isb.edu.pk/VerifyLogin.asp?sid=974494673"

	data := url.Values{}
	data.Set("txtLoginName", username)
	data.Set("txtPassword", password)
	data.Set("txtCampus_Id", "1")

	resp, err := client.PostForm(loginURL, data)
	if err != nil {
		fmt.Printf("Request Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read Error: %v\n", err)
		return
	}

	fmt.Printf("Final Destination: %s\n", resp.Request.URL.String())
	fmt.Println("--- Page Content Start ---")
	fmt.Println(string(body))
	fmt.Println("--- Page Content End ---")
}
