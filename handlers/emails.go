package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"

	"solar-garlic-mailing-list/database"
	"solar-garlic-mailing-list/email"
	"solar-garlic-mailing-list/model"

	"github.com/go-playground/validator/v10"
)

// really a factory that creates a handler function
func CreateEmailHandler(db *sql.DB, validate *validator.Validate, mailConfig email.MailConfig, skip_verify, verifyEndpoint string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the email from the request body

		var body = &model.CreateEmailBody{}

		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			slog.Info("CreateEmailHandler received invalid request body", "body", *body)
			errorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		emailAddress := body.Email

		if errs := validate.Var(emailAddress, "required,email"); errs != nil {
			slog.Error("Invalid email address", "email", emailAddress, "errs", errs)
			errorResponse(w, "Invalid body parameter: email. Please use a valid email address.", http.StatusBadRequest)
			return
		}

		// token for email verification
		token, err := email.CreateVerificationCode()
		if err != nil {
			slog.Error("Could not generate email verification token", "email", emailAddress)
			errorResponse(w, "Error creating user with that email", http.StatusInternalServerError)
		}

		unsubscribeID, err := email.CreateUnsubscribeId()
		if err != nil {
			slog.Error("Could not generate email unsubscribe id", "email", emailAddress)
			errorResponse(w, "Error creating user with that email", http.StatusInternalServerError)
		}

		e, t, dbErr := database.CreateEmail(db, emailAddress, token, unsubscribeID)

		if dbErr != nil {
			slog.Error("Database error while creating email", "err", dbErr)

			// check if the email already exists
			existing, err := database.ListEmails(db, emailAddress)
			if err != nil {
				slog.Error("Database error while creating email", "err", err)
				errorResponse(w, "Error creating user with that email", http.StatusInternalServerError)
				return
			}

			if len(existing) > 0 {
				slog.Error("Tried to create email address that already exists", "email", e)
				errorResponse(w, "Email address already exists", http.StatusConflict)
				return
			}

			errorResponse(w, "Error creating user with that email", http.StatusInternalServerError)
			return
		}

		if skip_verify != "true" {
			// send a verification email in the background
			go email.SendVerificationEmail(
				mailConfig.MailFrom,
				emailAddress,
				mailConfig.User,
				mailConfig.Password,
				mailConfig.Host,
				t.Token,
				verifyEndpoint,
			)
		}

		w.Header().Set("Content-Type", "application/json")

		encodeErr := json.NewEncoder(w).Encode(e)
		if encodeErr != nil {
			slog.Error("Error encoding CreateEmailHandler response body", "err", encodeErr)
			errorResponse(w, "Error encoding response", http.StatusInternalServerError)
			return
		}

	}
}

func ListEmailsHandler(db *sql.DB, validate *validator.Validate) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		email, decodeErr := decodeEmailAddressParam(r)
		if decodeErr != nil {
			slog.Error("Error while decoding query string param", "param", "a", "err", decodeErr)
			errorResponse(w, "Invalid query string param: a. Please use a url-safe base64 encoded email.", http.StatusBadRequest)
			return
		}
		if email != "" {
			// TODO: validate the email
			slog.Info("Querying for emails by address.", "email", email)
			if errs := validate.Var(email, "required,email"); errs != nil {
				slog.Error("Validation error for email", "email", email, "errs", errs)
				errorResponse(w, "Invalid query string param: a. Please use a url-safe base64 encoded email.", http.StatusBadRequest)
				return
			}
		}
		emails, dbErr := database.ListEmails(db, email)
		if dbErr != nil {
			slog.Error("Database error while listing emails", "err", dbErr)
			errorResponse(w, "Error listings emails", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		encodeErr := json.NewEncoder(w).Encode(model.ListEmailsResponse{
			Emails: emails,
		})
		if encodeErr != nil {
			slog.Error("Error encoding ListEmailsHandler response body", "err", encodeErr)
			errorResponse(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

func VerifyEmail(db *sql.DB, validate *validator.Validate, config email.MailConfig, ttlSeconds int, redirectURL, unsubscribeEndpoint string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		token := params.Get("t")

		// TODO: validate token format
		if errs := validate.Var(token, `required,base64rawurl`); errs != nil {
			slog.Error("Invalid email verification token", "token", token, "errs", errs)
			errorResponse(w, "Invalid url parameter: t. Please use a valid verification token.", http.StatusBadRequest)
			return
		}

		slog.Info("Verifying mailing list registration token", "t", token)

		e, unsubscribeID, err := database.VerifyToken(db, token, "email", ttlSeconds)

		if err != nil {
			slog.Error("Could not verify mailing list token", "token", token, "err", err)
			errorResponse(w, "Could not verify mailing list token", http.StatusInternalServerError)

			// TODO: redirect the client to an error page instead
			return
		}

		// send sign up success email
		slog.Info("Sending mailing list sign-up success email", "email", e)
		go email.SendSignUpSuccessEmail(config.MailFrom, e, config.User, config.Password, config.Host, unsubscribeID, unsubscribeEndpoint)

		// TODO: redirect the client to a success page
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}

func UnsubscribeEmail(db *sql.DB, validate *validator.Validate) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		unsubscribeId := params.Get("id")

		// TODO: validate token format

		slog.Info("Unsubscribing mailing list user", "unsubscribeId", unsubscribeId)

		err := database.UnsubscribeEmail(db, unsubscribeId)
		if err != nil {
			slog.Error("Could not unsubscribe mailing list user", "unsubscribeId", unsubscribeId, "err", err)
			errorResponse(w, "Could not unsubscribe mailing list user", http.StatusInternalServerError)

			return
		}

		slog.Info("Successfully unsubscribed mailing list user", "unsubscribeId", unsubscribeId)

		w.Header().Set("Content-Type", "application/json")
		encodeErr := json.NewEncoder(w).Encode(model.UnsubscribeResponse{
			ID: unsubscribeId,
			Message: "Successfully unsubscribed user from mailing list",
		})
		if encodeErr != nil {
			slog.Error("Error encoding ListEmailsHandler response body", "err", encodeErr)
			errorResponse(w, "Error encoding response", http.StatusInternalServerError)
			return
		}		
	}
}

func decodeEmailAddressParam(r *http.Request) (string, error) {
	query := r.URL.Query()

	encodedEmail := query.Get("a")
	email, err := base64.RawURLEncoding.DecodeString(encodedEmail)

	if err != nil {
		return "", err
	}

	return string(email), nil
}

func errorResponse(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(model.ErrorResponse{Message: msg})
	if err != nil {
		slog.Error("Could not encode error response", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
