package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

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
