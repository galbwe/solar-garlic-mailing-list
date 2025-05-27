package email

import (
	"fmt"
	"log/slog"
	"solar-garlic-mailing-list/utils"

	"gopkg.in/mail.v2"
)

func SendVerificationEmail(from, to, user, password, host, token string) {
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Solar Garlic Band - Email Verification")
	// TODO: move verification host to environment variables
	verificationLink := fmt.Sprintf("http://localhost:8080/api/v1/emails/verify?t=%v", token)
	body := `
	<p>Let's Get Funky!</p>
	<p>But first, click <a href=%q>here</a> to finish signing up for Solar Garlic's mailing list. You'll get the latest news on our upcoming shows and releases %v</p>
	<p>Jam On %v</p>
	`
	m.SetBody("text/html", fmt.Sprintf(body, verificationLink, "😎", "🤘🏽"))

	// TODO: refactor duplicated email sending into a helper function
	d := mail.NewDialer(
		host,
		587,
		user,
		password,
	)
	if err := d.DialAndSend(m); err != nil {
		slog.Error("Failed to send email", "err", err)
	}

	slog.Info("Sent verification email", "to", to)
}

func SendSignUpSuccessEmail(from, to, user, password, host, unsubscribeID string) {
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Solar Garlic Band - You're Subscribed!")

	// TODO: move host for unsubscribe link to environment variables
	unsubscribeLink := fmt.Sprintf("<http://localhost:8080/api/v1/emails/unsubscribe?id=%v>", unsubscribeID)
	m.SetHeader("List-Unsubscribe", unsubscribeLink)

	body := `
	<p>You're Subscribed!</p>	
	<p>Get ready to receive the latest updates on Solar Garlic's shows and releases.</p>
	<p>It's gonna be awesome! 🎸</p>
	`
	m.SetBody("text/html", body)

	d := mail.NewDialer(
		host,
		587,
		user,
		password,
	)
	if err := d.DialAndSend(m); err != nil {
		slog.Error("Failed to send email", "err", err)
	}

	slog.Info("Sent sign up success email", "to", to)

}

func CreateVerificationCode() (string, error) {
	token, err := utils.GenerateToken(20)
	if err != nil {
		return "", err
	}
	return token, nil
}

func CreateUnsubscribeId() (string, error) {
	token, err := utils.GenerateToken(20)
	if err != nil {
		return "", err
	}
	return token, nil
}
