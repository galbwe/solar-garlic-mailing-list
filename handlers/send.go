package handlers

import "net/http"

type Send struct{}

func (h *Send) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sending emails ...\n"))
}
