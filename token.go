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
	"strings"
	"time"
)

func getTokenFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, tokenDir, tokenFile)
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
