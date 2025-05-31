package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(nt.Time.Format(time.RFC3339))
}

type Email struct {
	ID            int       `json:"id"`
	Email         string    `json:"email"`
	DateCreated   time.Time `json:"dateCreated"`
	DateVerified  NullTime  `json:"dateVerified"`
	Subscribed    bool      `json:"subscribed"`
	UnsubscribeId string    `json:"-"`
}

type Token struct {
	ID          int       `json:"id"`
	EmailID     int       `json:"emailID"`
	Token       string    `json:"token"`
	Type        string    `json:"type"`
	DateCreated time.Time `json:"dateCreated"`
}

type CreateEmailBody struct {
	Email string `json:"email"`
}

type ListEmailsResponse struct {
	Emails []Email `json:"emails"`
}

type UnsubscribeResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"msg"`
}
