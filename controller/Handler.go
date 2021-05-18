package controller

import (
	"net/http"

	"github.com/go-chi/chi"
	"Friend_management/db"
)

var DBInstance db.Database

func NewHandler(db db.Database) http.Handler {
	router := chi.NewRouter()	
	DBInstance = db
	router.Route("/users", Users)
	router.Route("/relationship", Relationship)
	return router
}
