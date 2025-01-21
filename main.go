package main

import (
	"net/http"
	"solar-garlic-mailing-list/handlers"
)

func main() {
	email := http.NewServeMux()
	email.Handle("/subscribe", &handlers.Subscribe{})
	email.Handle("/verify", &handlers.Verify{})
	email.Handle("/unsubscribe", &handlers.Unsubscribe{})

	admin := http.NewServeMux()
	admin.Handle("/send", &handlers.Send{})
	admin.Handle("/schedule", &handlers.Schedule{})

	// root muxer
	mux := http.NewServeMux()
	mux.Handle("/email/", http.StripPrefix("/email", email))
	mux.Handle("/admin/", http.StripPrefix("/admin", admin))
	mux.Handle("/healthcheck", &handlers.Healthcheck{})

	http.ListenAndServe(":8080", mux)
}
