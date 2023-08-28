package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func listCards(token *TokenResponse, deckID int) []CardResponse {
	defer waitForAnyKey()
	err := checkAndRenewToken(token)
	if err != nil {
		printError("Error renewing token")
		return nil
	}

	url := baseURL + "/card/getAll"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	query := req.URL.Query()
	query.Add("deck_id", strconv.Itoa(deckID))
	req.URL.RawQuery = query.Encode()

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		printError(err.Error())
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var cards []CardResponse
		err := json.NewDecoder(resp.Body).Decode(&cards)
		if err != nil {
			printError(err.Error())
			return nil
		}
		printSuccess("Success fetching the cards")
		return cards
	}

	body, _ := ioutil.ReadAll(resp.Body)
	printError("Fetching cards failed due to " + string(body))
	return nil
}

func addCard(token *TokenResponse, deckID int) {
	defer waitForAnyKey()
	err := checkAndRenewToken(token)
	if err != nil {
		printError("Error renewing token")
		return
	}

	question, answer := addCardForm()
	if question == "" || answer == "" {
		return
	}
	payload := url.Values{}
	payload.Set("deck_id", strconv.Itoa(deckID))
	payload.Set("question", question)
	payload.Set("answer", answer)
	payloadStr := payload.Encode()

	url := baseURL + "/card/create"
	req, err := http.NewRequest("POST", url, strings.NewReader(payloadStr))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		printError(err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		printSuccess("Success adding a card")
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	printError("Adding card failed due to " + string(body))
}
