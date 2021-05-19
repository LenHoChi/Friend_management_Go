package handler
import (
	"net/http"

	"github.com/go-chi/chi"
	"Friend_management/db"
	"Friend_management/util"
	"Friend_management/controller"
)

func NewHandler(db db.Database) http.Handler {
	router := chi.NewRouter()	
	util.DBInstance = db
	router.Route("/users", controller.Users)
	router.Route("/relationship", controller.Relationship)
	return router
}
