package main

import (
	"net/http"
)

func main() {
	email := http.NewServeMux()
	email.Handle("/subscribe", &subscribeHandler{})
	email.Handle("/verify", &verifyHandler{})
	email.Handle("/unsubscribe", &unsubscribeHandler{})

	admin := http.NewServeMux()
	admin.Handle("/send", &adminSendHandler{})
	admin.Handle("/schedule", &adminScheduleHandler{})

	// root muxer
	mux := http.NewServeMux()
	mux.Handle("/", &healthcheckHandler{})
	mux.Handle("/healthcheck", &healthcheckHandler{})
	mux.Handle("/email/", http.StripPrefix("/email", email))
	mux.Handle("/admin/", http.StripPrefix("/admin", admin))

	http.ListenAndServe(":8080", mux)
}

type healthcheckHandler struct{}

func (h *healthcheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy!"))
}

type  subscribeHandler struct{}

func (h *subscribeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Handling Email Registration ...\n"))
}

type verifyHandler struct{}

func (h *verifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Handling Email Verification ...\n"))
}

type adminSendHandler struct{}

func (h *adminSendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sending Message ...\n"))
}

type adminScheduleHandler struct{}

func (h *adminScheduleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Scheduling Message For Later ...\n"))
}

type  unsubscribeHandler struct{}

func (h *unsubscribeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Unsubscribing ...\n"))
}
