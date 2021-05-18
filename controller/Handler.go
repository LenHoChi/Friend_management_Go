package controller

import (
	// "database/sql"
	"net/http"

	"github.com/go-chi/chi"
	// "github.com/go-chi/render"
	"Friend_management/db"
)

// var DB2 db.Database2
var DBInstance db.Database
// var DB3 *sql.Tx
// func NewHandler(db db.Database2) http.Handler {
// 	router := chi.NewRouter()
// 	DB2 = db

// 	router.Route("/users", Users)
// 	router.Route("/relationship", Relationship)
// 	return router
// }
// func NewHandler2(db *sql.Tx) http.Handler {
// 	router := chi.NewRouter()	
// 	DB3 = db
// 	router.Route("/users", Users)
// 	router.Route("/relationship", Relationship)
// 	return router
// }
func NewHandler(db db.Database) http.Handler {
	router := chi.NewRouter()	
	DBInstance = db
	router.Route("/users", Users)
	router.Route("/relationship", Relationship)
	return router
}
