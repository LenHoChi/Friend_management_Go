package controller

import (
	"fmt"
	"net/http"

	r_Response "Friend_management/models/response"
	

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"Friend_management/db"
	"Friend_management/models"
	repo "Friend_management/repository"
	ser "Friend_management/services"
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
	users, err := userServices.FindAllUser(DBInstance)
	if err != nil {
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError,err.Error())
		return
	}
	if err := render.Render(w, r, users); err != nil {
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError,err.Error())
		return 
	}
	// r_Response.ResponseWithJSON(w, http.StatusOK, "ok")	
}

func (*controller)CreateUser(w http.ResponseWriter, r *http.Request){
	x,_ := db.Conn.Begin()
	DBInstance.Conn = x
	user := &models.User{}
	if err := render.Bind(r, user); err != nil {
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
		fmt.Println(err.Error())
		return
	}
	if !repo.IsEmailValid(user.Email){
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, "email was wrong")
		return
	}
	fmt.Println("d",user.Email)
	if err:= userServices.AddUser(DBInstance, user); err != nil {
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
		return 
	}
	if err := render.Render(w, r, user); err != nil {
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
		return     
	}
	DBInstance.Conn.Commit()
}

func (*controller)GetUser(w http.ResponseWriter, r *http.Request) {
	// email := r.Context().Value("emailKey").(string)
	email := r.URL.Query().Get("id")
	if !repo.IsEmailValid(email){
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, "email was wrong")
		return
	}
	user, err := userServices.GetUserByEmail(DBInstance, email)
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
	DBInstance.Conn = x
	email := r.URL.Query().Get("id")
	if !repo.IsEmailValid(email){
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, "email was wrong")
		return
	}
	err := userServices.DeleteUser(DBInstance, email)
	if err != nil {
		if err == db.ErrNoMatch {
			r_Response.ResponseWithJSON(w, http.StatusInternalServerError, "No any user match")
		}else{
			r_Response.ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	DBInstance.Conn.Commit()
}
