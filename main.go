package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"solar-garlic-mailing-list/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var db *sql.DB

func main() {

	db = database.CreateDB()
	database.CreateSchema(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/emails", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("list emails"))
	})
	r.Post("/emails", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("create email"))
	})

	r.Get("/emails/{id:^[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		w.Write([]byte(fmt.Sprintf("get email with id %q", id)))
	})
	r.Put("/emails/{id:^[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		w.Write([]byte(fmt.Sprintf("update email with id %q", id)))
	})
	r.Delete("/emails/{id:^[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		w.Write([]byte(fmt.Sprintf("delete email with id %q", id)))
	})

	http.ListenAndServe(":8080", r)

}
