package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"solar-garlic-mailing-list/model"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB(file string) *sql.DB {
	dsn := fmt.Sprintf("file:%v?cache=shared", file)
	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		panic(fmt.Sprintf("Could not get a database connection: %v\n", err))
	}

	return db
}

func CreateSchema(db *sql.DB) {
	// create the emails table
	stmt := `
		CREATE TABLE IF NOT EXISTS emails (
			id INTEGER NOT NULL PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
			date_verified DATETIME DEFAULT NULL,
			subscribed BOOLEAN DEFAULT true
		);
	`
	slog.Info("Creating emails table", "stmt", stmt)
	_, err := db.Exec(stmt)
	if err != nil {
		slog.Error("Error creating emails table", "err", err)
		panic(fmt.Sprintf("Could not create emails table: %v\n", err))
	}

	// create the tokens table
	stmt = `
		CREATE TABLE IF NOT EXISTS tokens (
			id INTEGER NOT NULL PRIMARY KEY,
			email_id INTEGER,	
			token TEXT NOT NULL UNIQUE,
			type TEXT NOT NULL CHECK (type in ('email', 'auth')),
			date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(email_id) REFERENCES emails(id)
		);
	`
	slog.Info("Creating the tokens table", "stmt", stmt)
	_, err = db.Exec(stmt)
	if err != nil {
		slog.Error("Error creating tokens table", "err", err)
		panic(fmt.Sprintf("Could not create tokens table: %v\n", err))
	}
}

func CreateEmail(db *sql.DB, email string, token string) (*model.Email, *model.Token, error) {
	var createdEmail = &model.Email{}
	var createdToken = &model.Token{}

	tx, err := db.Begin()

	if err != nil {
		slog.Error("Error creating database transaction in CreateEmail", "err", err)
		return nil, nil, err
	}

	// Create the email
	stmt := `
		INSERT INTO emails (email) VALUES ($1)
		RETURNING id, email, date_created, date_verified, subscribed
	`
	params := []any{email}
	slog.Info("Inserting email into database", "stmt", stmt, "params", params)

	err = tx.QueryRow(stmt, params...).Scan(
		&createdEmail.ID,
		&createdEmail.Email,
		&createdEmail.DateCreated,
		&createdEmail.DateVerified,
		&createdEmail.Subscribed,
	)

	if err != nil {
		slog.Error("Error inserting email into database", "stmt", stmt, "params", params, "err", err)
		tx.Rollback()
		return nil, nil, err
	}

	// Create a verification token for the email
	stmt = `
		INSERT INTO tokens (email_id, token, type) 
		VALUES ($1, $2, 'email')
		RETURNING id, email_id, token, type, date_created
	`
	params = []any{createdEmail.ID, token}
	slog.Info("Inserting token into database", "stmt", stmt, "params", params)
	err = tx.QueryRow(stmt, params...).Scan(
		&createdToken.ID,
		&createdToken.EmailID,
		&createdToken.Token,
		&createdToken.Type,
		&createdToken.DateCreated,
	)
	if err != nil {
		slog.Error("Error inserting token into database", "stmt", stmt, "params", params, "err", err)
	}
	tx.Commit()
	return createdEmail, createdToken, nil
}

func ListEmails(db *sql.DB, email string) ([]model.Email, error) {

	stmt, params := getListEmailsQuery(email)

	slog.Info("Selecting emails from database", "stmt", stmt, "params", params)

	rows, err := db.Query(stmt, params...)
	if err != nil {
		slog.Error("Error selecting email list from database", "stmt", stmt, "params", params)
		return nil, err
	}
	defer rows.Close()

	emails := []model.Email{}

	for rows.Next() {
		var email model.Email

		scanErr := rows.Scan(
			&email.ID,
			&email.Email,
			&email.DateCreated,
			&email.DateVerified,
			&email.Subscribed,
		)

		if scanErr != nil {
			slog.Error("Error scanning sql rows", "stmt", stmt, "params", params, "err", scanErr)
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}

func getListEmailsQuery(email string) (string, []any) {
	if email == "" {
		stmt := `
			SELECT  
				id,
				email,
				date_created,
				date_verified,
				subscribed	
			FROM emails;	
		`
		params := []any{}
		return stmt, params
	}

	stmt := `
		SELECT	
			id,			
			email,
			date_created,
			date_verified,
			subscribed
		FROM emails
		where email = $1;
	`
	params := []any{email}
	return stmt, params
}

func VerifyToken(db *sql.DB, token, tokenType string, ttlSeconds int) error {
	// get emails from the database with a matching token that has not expired yet
	var emailId int
	stmt := `
		SELECT e.id 
		FROM emails e join tokens t on e.id = t.email_id
		WHERE
			t.token = $1
			AND t.type = $2
			AND CURRENT_TIMESTAMP < datetime(t.date_created, $3)
			ORDER BY t.date_created DESC
			LIMIT 1;
	`
	params := []any{token, tokenType, fmt.Sprintf(`+%v seconds`, ttlSeconds)}

	err := db.QueryRow(stmt, params...).Scan(&emailId)

	if err != nil {
		slog.Error("Error querying database for emails", "stmt", stmt, "params", params, "err", err)
		return err
	}

	// return an error if no email with a matching token is found

	// update the email verified date
	if tokenType == "email" {
		stmt = `
			UPDATE emails
			SET date_verified = CURRENT_TIMESTAMP	
			WHERE id = $1;
		`
		params = []any{emailId}

		_, err = db.Exec(stmt, params...)

		if err != nil {
			slog.Error("Error updating email in database", "stmt", stmt, "params", params, "err", err)
			return err
		}
	}

	// TODO: add verfication for admin user auth tokens

	return nil
}
