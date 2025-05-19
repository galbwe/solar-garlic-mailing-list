package model

import (
	"database/sql"
	"time"
)

type Email struct {
	ID           int          `json:"id"`
	Email        string       `json:"email"`
	DateCreated  time.Time    `json:"dateCreated"`
	DateVerified sql.NullTime `json:"dateVerified"`
	Subscribed   bool         `json:"subscribed"`
}

type CreateEmailBody struct {
	Email string `json"email"`
}

type ListEmailsResponse struct {
	Emails []Email `json:"emails"`
}
