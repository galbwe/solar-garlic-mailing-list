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
	verificationLink := fmt.Sprintf("http://localhost:8080/api/v1/emails/verify?t=%v", token)
	body := `
	<p>Let's Get Funky!</p>
	<p>But first, click <a href=%q>here</a> to finish signing up for Solar Garlic's mailing list. You'll get the latest news on our upcoming shows and releases %v</p>
	<p>Jam On %v</p>
	`
	m.SetBody("text/html", fmt.Sprintf(body, verificationLink, "😎", "🤘🏽"))

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

func CreateVerificationCode() (string, error) {
	token, err := utils.GenerateToken(20)
	if err != nil {
		return "", err
	}
	return token, nil
}
