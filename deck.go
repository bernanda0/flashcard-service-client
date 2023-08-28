package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func listDecks(token *TokenResponse, accountID int) []DeckResponse {
	defer waitForAnyKey()
	err := checkAndRenewToken(token)
	if err != nil {
		printError("Error renewing token")
		return nil
	}

	url := baseURL + "/deck/getAll"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	query := req.URL.Query()
	query.Add("account_id", strconv.Itoa(accountID))
	req.URL.RawQuery = query.Encode()

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		printError(err.Error())
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var decks []DeckResponse
		err := json.NewDecoder(resp.Body).Decode(&decks)
		if err != nil {
			printError(err.Error())
			return nil
		}
		printSuccess("Success fetching the decks")
		return decks
	}

	body, _ := ioutil.ReadAll(resp.Body)
	printError("Fetching decks failed due to " + string(body))
	return nil
}

func addDeck(token *TokenResponse, accountID int) {
	// Use the token to authenticate requests
	// Send a POST request to add a new card
	// Get user input for card details
	// Handle the response
}
