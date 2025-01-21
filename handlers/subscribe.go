package handlers

import "net/http"

type  Subscribe struct{}

func (h *Subscribe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Handling Email Registration ...\n"))
	// get user's email from json request
	
}