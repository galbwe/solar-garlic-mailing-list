package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"solar-garlic-mailing-list/helpers"
)

type SubscribeBody struct {
	Email string
}

type Subscribe struct{}

func (h *Subscribe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Handling Email Registration ...\n"))
	// get user's email from json request
	var body SubscribeBody
	err := helpers.DecodeJSONBody(w, r, &body)
	if err != nil {
		var mr *helpers.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	fmt.Fprintf(w, "body: %+v", body)
	// TODO: validate the email format

	// TODO: send a verification email

}
