package handlers

import "net/http"

type Verify struct{}

func (h *Verify) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Verifying email address with token ...\n"))
}
