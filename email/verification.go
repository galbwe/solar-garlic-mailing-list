package email

import (
	"fmt"
	"log/slog"

	"gopkg.in/mail.v2"
)

func SendVerificationEmail(from, to, user, password, host string) {
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Solar Garlic Band - Email Verification")
	verificationLink := `https://en.wikipedia.org/wiki/Mailing_list`
	body := `
	<p>Let's Get Funky!</p>
	<p>But first, click the link below to finish signing up for Solar Garlic's mailing list. You'll get the latest news on our upcoming shows and releases %v</p>
	<a href=%q>%v</a>
	<p>Jam On %v</p>
	`
	m.SetBody("text/html", fmt.Sprintf(body, "😎", verificationLink, verificationLink, "🤘🏽"))

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
