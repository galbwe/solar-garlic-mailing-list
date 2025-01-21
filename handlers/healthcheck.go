package handlers

import "net/http"

type Healthcheck struct{}

func (h *Healthcheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy!"))
}
