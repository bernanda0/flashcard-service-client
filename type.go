package main

import "time"

type TokenResponse struct {
	SessionID          string    `json:"session_id"`
	AccessToken        string    `json:"access_token"`
	AccessTokenExpire  time.Time `json:"access_token_expire"`
	RefreshToken       string    `json:"refresh_token"`
	RefreshTokenExpire time.Time `json:"refresh_token_expire"`
	UserID             int       `json:"user_id"`
	Username           string    `json:"username"`
}
