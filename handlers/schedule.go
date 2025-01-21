package handlers

import "net/http"

type Schedule struct{}

func (h *Schedule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Scheduling mailing list blast ...\n"))
	// get user's email from json request

}
