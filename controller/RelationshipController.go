package controller

import (
	// "context"
	// "fmt"
	// "errors"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	// "Friend_management/db"
	// "Friend_management/models"
	r_Request "Friend_management/models/request"
	// "encoding/json"
	// "github.com/gin-gonic/gin"
	"Friend_management/db"
	"Friend_management/exception"
	repo "Friend_management/repository"
	ser "Friend_management/services"
)

var RelationshipKey = "relationshipKey"
var (
	repoS repo.RelationshipInter = repo.NewRepoRelationship()
	serS ser.RepositoryService = ser.NewRelationshipService(repoS) 
)
func Relationship (router chi.Router) {
	router.Get("/", NewRelationshipControl(serS).GetAllRelationships)
	router.Post("/make", NewRelationshipControl(serS).MakeFriend)
	router.Post("/list", NewRelationshipControl(serS).FindListFriend)
	router.Post("/common", NewRelationshipControl(serS).FindCommonListFriend)
	router.Post("/update",NewRelationshipControl(serS).BeSubcriber)
	router.Post("/block", NewRelationshipControl(serS).ToBLock)
	router.Post("/retrieve", NewRelationshipControl(serS).RetrieveUpdate)
}
var (
	relationshipService ser.RepositoryService
)
type controllerRelationship struct{}
type RelationshipControl interface{
	GetAllRelationships (w http.ResponseWriter, r *http.Request)
	MakeFriend(w http.ResponseWriter, r *http.Request)
	FindListFriend(w http.ResponseWriter, r *http.Request)
	FindCommonListFriend(w http.ResponseWriter, r *http.Request)
	BeSubcriber(w http.ResponseWriter, r *http.Request)
	ToBLock(w http.ResponseWriter, r *http.Request)
	RetrieveUpdate(w http.ResponseWriter, r *http.Request)
}
func NewRelationshipControl(ser ser.RepositoryService) RelationshipControl{
	relationshipService = ser
	return &controllerRelationship{}
}
func (*controllerRelationship)GetAllRelationships (w http.ResponseWriter, r *http.Request) {
	relationships, err := relationshipService.GetAllRelationship(DBInstance)
	fmt.Println("loi:",err)
	if err != nil {
		// responseWithError(w, http.StatusBadRequest ,err.Error())
		render.Render(w,r,exception.ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, relationships); err != nil {
		render.Render(w,r, exception.ErrorRenderer(err))
		return
	}
	//DBInstance.Conn.Commit()
	// DBInstance.Conn.Rollback()
}
// {"friends":["1","2"]}
func (*controllerRelationship)MakeFriend(w http.ResponseWriter, r *http.Request){
	x,_ := db.Conn.Begin()
	DBInstance.Conn = x
	requestAddFriend := &r_Request.RequestFriendLists{}
	render.Bind(r, requestAddFriend)
	//check length <2
	userEmail := requestAddFriend.RequestFriendLists[0]
	friendEmail := requestAddFriend.RequestFriendLists[1]
	//check valid for two emails
	if !isEmailValid(userEmail)||!isEmailValid(friendEmail){
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("email is wrong")))
		return
	}
	responseRS, err := relationshipService.AddRelationship(DBInstance, userEmail, friendEmail)
	if err != nil {
		render.Render(w,r,exception.ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, responseRS); err != nil {
		render.Render(w,r, exception.ErrorRenderer(err))
		return
	}
	DBInstance.Conn.Commit()
}
//{"email":"1"}
func (*controllerRelationship)FindListFriend(w http.ResponseWriter, r *http.Request){
	Argument := &r_Request.RequestEmail{}
	render.Bind(r, Argument)
	if !isEmailValid(Argument.Email){
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("email is wrong")))
		return
	}
	responseRS, err := relationshipService.FindListFriend(DBInstance, Argument.Email)
	if err != nil{
		render.Render(w,r,exception.ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, responseRS); err != nil {
		render.Render(w,r, exception.ErrorRenderer(err))
		return
	}
}
// {"friends":["1","2"]}
func (*controllerRelationship)FindCommonListFriend(w http.ResponseWriter, r *http.Request){
	rsFriend := &r_Request.RequestFriendLists{}
	ls := make([]string,0)
	render.Bind(r, rsFriend)
	if !isEmailValid(rsFriend.RequestFriendLists[0])||!isEmailValid(rsFriend.RequestFriendLists[1]){
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("email is wrong")))
		return
	}
	ls = append(ls, rsFriend.RequestFriendLists[0], rsFriend.RequestFriendLists[1])
	lst, err := relationshipService.FindCommonListFriend(DBInstance, ls)
	if err != nil {
		render.Render(w,r,exception.ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w,r,lst); err != nil {
		render.Render(w,r, exception.ErrorRenderer(err))
		return
	}
}
// {"requestor":"len1","target":"len2"}
func (*controllerRelationship)BeSubcriber(w http.ResponseWriter, r *http.Request){
	x,_ := db.Conn.Begin()
	DBInstance.Conn = x
	Argument := &r_Request.RequestUpdate{}
	render.Bind(r, Argument)
	if !isEmailValid(Argument.Requestor)||!isEmailValid(Argument.Target){
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("email is wrong")))
		return
	}
	responseRS, err:= relationshipService.BeSubcribe(DBInstance, Argument.Requestor, Argument.Target)
	if err != nil {
		render.Render(w,r,exception.ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, responseRS); err != nil {
		render.Render(w,r, exception.ErrorRenderer(err))
		return
	}
	DBInstance.Conn.Commit()
}
func (*controllerRelationship)ToBLock(w http.ResponseWriter, r *http.Request){
	x,_ := db.Conn.Begin()
	DBInstance.Conn = x
	Argument := &r_Request.RequestUpdate{}
	render.Bind(r, Argument)
	if !isEmailValid(Argument.Requestor)||!isEmailValid(Argument.Target){
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("email is wrong")))
		return
	}
	responseRS ,err :=relationshipService.ToBlock(DBInstance, Argument.Requestor, Argument.Target)
	if err != nil {
		render.Render(w,r,exception.ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, responseRS); err != nil {
		render.Render(w,r, exception.ErrorRenderer(err))
		return
	}
	DBInstance.Conn.Commit()
}
// {"sender":"len1","target":"len2"}
func (*controllerRelationship)RetrieveUpdate(w http.ResponseWriter, r *http.Request){
	Argument := &r_Request.RetrieveUpdate{}
	render.Bind(r, Argument)
	if !isEmailValid(Argument.Sender){
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("email is wrong")))
		return
	}
	responseRS, err := relationshipService.RetrieveUpdate(DBInstance,Argument.Sender, Argument.Tartget)
	if err != nil {
		render.Render(w,r,exception.ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, responseRS); err != nil {
		render.Render(w,r, exception.ErrorRenderer(err))
		return
	}
}


