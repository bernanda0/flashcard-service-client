package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	baseURL   = "http://localhost:4444" // Replace with your server's URL
	tokenFile = "token.json"            // Changed to token.json
)

type TokenResponse struct {
	SessionID          string    `json:"session_id"`
	AccessToken        string    `json:"access_token"`
	AccessTokenExpire  time.Time `json:"access_token_expire"`
	RefreshToken       string    `json:"refresh_token"`
	RefreshTokenExpire time.Time `json:"refresh_token_expire"`
	UserID             int       `json:"user_id"`
	Username           string    `json:"username"`
}

func main() {
	for {
		fmt.Println("Welcome to the flash card app!")

		token, err := readTokenFromFile()
		if err != nil {
			fmt.Println("[You're not authenticated. Please login]")
			fmt.Println("1. Login")
			fmt.Println("2. Exit")

			var choice int
			fmt.Print("Enter your choice: ")
			_, err := fmt.Scanln(&choice)
			if err != nil {
				fmt.Println("Invalid input. Please enter a number.")
				continue
			}

			switch choice {
			case 1:
				loginResult := login()
				if loginResult {
					continue // Successful login, go to the main menu
				}
				// Failed login, continue the loop

			case 2:
				fmt.Println("Goodbye!")
				return
			default:
				fmt.Println("Invalid choice. Please select a valid option.")
			}
		}

		// Display menu and handle user actions
		displayMenu()
		var choice int
		fmt.Print("Enter your choice: ")
		_, err = fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			listCards(token)
		case 2:
			addCard(token)
		case 3:
			logout()
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please select a valid option.")
		}
	}
}

func login() bool {
	var email, password string
	fmt.Print("Email: ")
	fmt.Scanln(&email)
	fmt.Print("Password: ")
	fmt.Scanln(&password)

	// Create a URL-encoded payload
	payload := url.Values{}
	payload.Set("email", email)
	payload.Set("password", password)
	payloadStr := payload.Encode()

	// Send the authentication request
	resp, err := http.Post(baseURL+"/auth/login", "application/x-www-form-urlencoded", strings.NewReader(payloadStr))
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {
		var tokenResponse TokenResponse
		err := json.NewDecoder(resp.Body).Decode(&tokenResponse)
		if err != nil {
			fmt.Println("Error decoding response:", err)
			return false
		}

		if tokenResponse.AccessToken != "" {
			saveTokenToFile(tokenResponse)
			fmt.Println("Login successful.")
			return true
		}
	} else {
		fmt.Println("Login failed.")
		return false
	}
	return false
}

// ...

func readTokenFromFile() (string, error) {
	tokenData, err := ioutil.ReadFile(tokenFile)
	if err != nil || len(tokenData) == 0 {
		return "", errors.New("no token found")
	}

	var token TokenResponse
	err = json.Unmarshal(tokenData, &token)
	if err != nil {
		return "", errors.New("invalid token format")
	}

	return token.AccessToken, nil
}

func saveTokenToFile(token TokenResponse) error {
	tokenData, err := json.Marshal(token)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(tokenFile, tokenData, 0644)
}

func displayMenu() {
	fmt.Println("Main Menu:")
	fmt.Println("1. My Cards")
	fmt.Println("2. Add Cards")
	fmt.Println("3. Logout")
}

func listCards(token string) {
	// Use the token to authenticate requests
	// Send a GET request to retrieve user's cards
	// Process the response and display the cards
}

func addCard(token string) {
	// Use the token to authenticate requests
	// Send a POST request to add a new card
	// Get user input for card details
	// Handle the response
}

func logout() {
	// Delete the token file
	_ = os.Remove(tokenFile)
	fmt.Println("Logged out.")
}
