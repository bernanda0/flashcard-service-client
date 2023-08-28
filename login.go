package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func login() {
	defer waitForAnyKey()
	email, password := loginForm()

	// Create a URL-encoded payload
	payload := url.Values{}
	payload.Set("email", email)
	payload.Set("password", password)
	payloadStr := payload.Encode()

	// Send the authentication request
	resp, err := http.Post(baseURL+"/auth/login", "application/x-www-form-urlencoded", strings.NewReader(payloadStr))
	if err != nil {
		printError(err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var tokenResponse TokenResponse
		err := json.NewDecoder(resp.Body).Decode(&tokenResponse)
		if err != nil {
			printError(err.Error())
			return
		}

		if tokenResponse.AccessToken != "" {
			err := saveTokenToFile(tokenResponse)
			if err == nil {
				printSuccess("Login success")
				return
			}
		}
	}
	body, _ := ioutil.ReadAll(resp.Body)
	printError("Login failed due to " + string(body))
}

func signup() bool {
	var username, email, password string
	fmt.Print("Username: ")
	fmt.Scanln(&username)
	fmt.Print("Email: ")
	fmt.Scanln(&email)
	fmt.Print("Password: ")
	fmt.Scanln(&password)

	// Create a URL-encoded payload
	payload := url.Values{}
	payload.Set("username", username)
	payload.Set("email", email)
	payload.Set("password", password)
	payloadStr := payload.Encode()

	// Send the authentication request
	resp, err := http.Post(baseURL+"/auth/signup", "application/x-www-form-urlencoded", strings.NewReader(payloadStr))
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Sign Up Success. Please")
		return true
	} else {
		fmt.Println("Sign Up failed.")
		return false
	}
}

func logout() {
	tokenPath := getTokenFilePath()

	err := os.Remove(tokenPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Logged out.")
}
