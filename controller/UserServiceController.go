package controller

import (
	// "context"
	// "fmt"
	// "errors"
	"fmt"
	"net/http"

	// "strconv"
	r_Response "Friend_management/models/response"
	"regexp"

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
	// router.Post("/", CreateUser)
	// router.Get("/find", GetUser)
	// router.Delete("/delete", DeleteUser)
	// router.Route("/{emailID}", func(router chi.Router) {
		// router.Use(UserContext)
		// router.Get("/", GetUser)
		// router.Delete("/", DeleteUser)
	// })
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
	if !isEmailValid(user.Email){
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
// func UserContext(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
// 		email := chi.URLParam(r, "emailID")
// 		if email == "" {
// 			r_Response.ResponseWithJSON(w, http.StatusInternalServerError, "Email is required")
// 			return
// 		}
// 		fmt.Println("day: ",email)
// 		ctx := context.WithValue(r.Context(), UserEmailKey, email)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func (*controller)GetUser(w http.ResponseWriter, r *http.Request) {
	// email := r.Context().Value("emailKey").(string)
	email := r.URL.Query().Get("id")
	if !isEmailValid(email){
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
	// email := r.Context().Value(UserEmailKey).(string)
	x,_ := db.Conn.Begin()
	DBInstance.Conn = x
	email := r.URL.Query().Get("id")
	if !isEmailValid(email){
		r_Response.ResponseWithJSON(w, http.StatusInternalServerError, "email was wrong")
		return
	}
	err := userServices.DeleteUser(DBInstance, email)
	if err != nil {
		fmt.Println("?")
		if err == db.ErrNoMatch {
			r_Response.ResponseWithJSON(w, http.StatusInternalServerError, "No any user match")
		}else{
			r_Response.ResponseWithJSON(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	DBInstance.Conn.Commit()
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !emailRegex.MatchString(e) {
		return false
	}
	return true
}