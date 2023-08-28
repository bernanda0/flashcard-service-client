package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL   = "http://localhost:4444"
	tokenFile = "token.json" // Changed to token.json
	tokenDir  = ".flashcardapp"
)

func getTokenFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, tokenDir, tokenFile)
}

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
					continue
				} else {
					continue
				}
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
			listDecks(token, token.UserID)
		case 2:
			addDeck(token.AccessToken, token.UserID)
		case 3:
			logout()
			fmt.Println("Goodbye!")
			return
		case 4:
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

func readTokenFromFile() (*TokenResponse, error) {
	tokenPath := getTokenFilePath()
	tokenData, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, errors.New("no token found")
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal(tokenData, &tokenResponse)
	if err != nil {
		return nil, errors.New("invalid token format")
	}

	return &tokenResponse, nil
}

func saveTokenToFile(tokenResponse TokenResponse) error {
	tokenPath := getTokenFilePath()
	tokenData, err := json.Marshal(tokenResponse)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(tokenPath), 0700)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(tokenPath, tokenData, 0600)
}

func checkAndRenewToken(tokenResponse *TokenResponse) error {
	if time.Now().After(tokenResponse.AccessTokenExpire) {
		req_url := baseURL + "/auth/renewToken"

		payload := url.Values{}
		payload.Set("refresh_token", tokenResponse.RefreshToken)
		fmt.Println(tokenResponse.RefreshToken)
		payloadStr := payload.Encode()

		resp, err := http.Post(req_url, "application/x-www-form-urlencoded", strings.NewReader(payloadStr))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusCreated {
			var renewedToken TokenResponse
			err := json.NewDecoder(resp.Body).Decode(&renewedToken)
			if err != nil {
				return err
			}

			tokenResponse.AccessToken = renewedToken.AccessToken
			tokenResponse.AccessTokenExpire = renewedToken.AccessTokenExpire

			// Save the updated token to token.json
			err = saveTokenToFile(*tokenResponse)
			if err != nil {
				return err
			}

			fmt.Println("Token renewed.")
		} else {
			return errors.New("token renewal failed with status: " + resp.Status)
		}
	}
	return nil
}

func displayMenu() {
	fmt.Println("Main Menu:")
	fmt.Println("1. My Deck")
	fmt.Println("2. Add Deck")
	fmt.Println("3. Logout")
	fmt.Println("4. Exit")
}

func listDecks(token *TokenResponse, accountID int) {
	err := checkAndRenewToken(token)
	if err != nil {
		fmt.Println("Error renewing token")
		return
	}

	url := baseURL + "/deck/getAll"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	query := req.URL.Query()
	query.Add("account_id", strconv.Itoa(accountID))
	req.URL.RawQuery = query.Encode()

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		fmt.Println(string(body))
	} else {
		fmt.Println("Request failed with status:", resp.Status)
	}
}

func addDeck(token string, account_id int) {
	// Use the token to authenticate requests
	// Send a POST request to add a new card
	// Get user input for card details
	// Handle the response
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
