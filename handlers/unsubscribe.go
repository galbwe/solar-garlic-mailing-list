package handlers

import "net/http"

type  Unsubscribe struct{}

func (h *Unsubscribe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Unsubscribing from mailing list ...\n"))
}