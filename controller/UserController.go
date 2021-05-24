package controller

import (
	"net/http"

	r_Response "Friend_management/models/response"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"Friend_management/db"
	"Friend_management/models"
	repo "Friend_management/repository"
	ser "Friend_management/services"
	"Friend_management/util"
)

var UserEmailKey = "emailKey"
var (
	userServices ser.UserService
)
type controller struct{}
type UserControl interface{
	GetAllUsers(response http.ResponseWriter, request *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}
func NewUserControl(ser ser.UserService) UserControl{
	userServices = ser
	return &controller{}
}
var (
	y repo.UserRepoInter = repo.NewRepo()
	x ser.UserService = ser.NewUserService(y)
)
func Users(router chi.Router) {
	router.Get("/", NewUserControl(x).GetAllUsers)
	router.Post("/",NewUserControl(x).CreateUser)
	router.Get("/find", NewUserControl(x).GetUser)
	router.Delete("/delete", NewUserControl(x).DeleteUser)
}
func (*controller)GetAllUsers(w http.ResponseWriter, r *http.Request) {
	x,_ := db.Conn.Begin()
	util.DBInstance.Conn = x
	users, err := userServices.FindAllUser(util.DBInstance)
	if err != nil {
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError,err.Error())
		return
	}
	if err := render.Render(w, r, users); err != nil {
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError,err.Error())
		return 
	}
}

func (*controller)CreateUser(w http.ResponseWriter, r *http.Request){
	x,_ := db.Conn.Begin()
	util.DBInstance.Conn = x
	user := &models.User{}
	if err := render.Bind(r, user); err != nil {
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !repo.IsEmailValid(user.Email){
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, "email was wrong")
		return
	}
	if err:= userServices.AddUser(util.DBInstance, user); err != nil {
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
		return 
	}
	if err := render.Render(w, r, user); err != nil {
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
		return     
	}
	util.DBInstance.Conn.Commit()
}

func (*controller)GetUser(w http.ResponseWriter, r *http.Request) {
	// email := r.Context().Value("emailKey").(string)
	x,_ := db.Conn.Begin()
	util.DBInstance.Conn = x
	email := r.URL.Query().Get("id")
	if len(email)==0{
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, "lack email")
		return
	}
	if !repo.IsEmailValid(email){
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, "email was wrong")
		return
	}
	user, err := userServices.GetUserByEmail(util.DBInstance, email)
	if err != nil {
		if err == db.ErrNoMatch {
			r_Response.ResponseWithJSON(w, http.StatusInternalServerError, "No any user match")
		}else {
			r_Response.ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if err := render.Render(w, r, &user); err != nil {
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
		return 
	}
}

func (*controller)DeleteUser(w http.ResponseWriter, r *http.Request) {
	x,_ := db.Conn.Begin()
	util.DBInstance.Conn = x
	email := r.URL.Query().Get("id")
	if len(email)==0{
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, "lack email")
		return
	}
	if !repo.IsEmailValid(email){
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, "email was wrong")
		return
	}
	err := userServices.DeleteUser(util.DBInstance, email)
	if err != nil {
		if err == db.ErrNoMatch {
			r_Response.ResponseWithJSON(w, http.StatusInternalServerError, "No any user match")
		}else{
			r_Response.ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	util.DBInstance.Conn.Commit()
}
