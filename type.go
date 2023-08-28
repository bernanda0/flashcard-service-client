package main

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
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

type DeckResponse struct {
	DeckID    int         `json:"deck_id"`
	AccountID int         `json:"account_id"`
	Title     string      `json:"title"`
	CreatedAt pq.NullTime `json:"created_at"`
}

type CardResponse struct {
	FlashcardID    int32           `json:"flashcard_id"`
	DeckID         int32           `json:"deck_id"`
	Question       string          `json:"question"`
	Answer         string          `json:"answer"`
	NextReviewDate sql.NullTime    `json:"next_review_date"`
	Interval       sql.NullInt32   `json:"interval"`
	Repetitions    sql.NullInt32   `json:"repetitions"`
	EasinessFactor sql.NullFloat64 `json:"easiness_factor"`
	CreatedAt      sql.NullTime    `json:"created_at"`
	UpdatedAt      sql.NullTime    `json:"updated_at"`
	IsArchived     sql.NullBool    `json:"is_archived"`
}
