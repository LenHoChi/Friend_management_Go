package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	r_Request "Friend_management/models/request"
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
		render.Render(w,r,exception.ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, relationships); err != nil {
		render.Render(w,r, exception.ErrorRenderer(err))
		return
	}
}
// {"friends":["1","2"]}
func (*controllerRelationship)MakeFriend(w http.ResponseWriter, r *http.Request){
	x,_ := db.Conn.Begin()
	DBInstance.Conn = x
	requestAddFriend := &r_Request.RequestFriendLists{}
	if err := render.Bind(r, requestAddFriend);err != nil{
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("invalid format")))
		return
	}
	//check length <2
	if len(requestAddFriend.RequestFriendLists)!=2{
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("not enough email")))
		return
	}
	userEmail := requestAddFriend.RequestFriendLists[0]
	friendEmail := requestAddFriend.RequestFriendLists[1]
	//check valid for two emails
	if !repo.IsEmailValid(userEmail)||!repo.IsEmailValid(friendEmail){
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("email is wrong")))
		return
	}
	if userEmail == friendEmail {
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("error cause 2 emails are same")))
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
	if err := render.Bind(r, Argument);err != nil{
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("invalid format")))
		return
	}
	if !repo.IsEmailValid(Argument.Email){
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
	if err := render.Bind(r, rsFriend);err != nil{
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("invalid format")))
		return
	}
	if len(rsFriend.RequestFriendLists)!=2{
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("enough email")))
		return
	}
	if !repo.IsEmailValid(rsFriend.RequestFriendLists[0])||!repo.IsEmailValid(rsFriend.RequestFriendLists[1]){
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("email is wrong")))
		return
	}
	if rsFriend.RequestFriendLists[0] == rsFriend.RequestFriendLists[1] {
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("error cause 2 emails are same")))
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
	if err := render.Bind(r, Argument);err != nil{
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("invalid format")))
		return
	}
	if !repo.IsEmailValid(Argument.Requestor)||!repo.IsEmailValid(Argument.Target){
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("email is wrong")))
		return
	}
	if Argument.Requestor == Argument.Target {
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("error cause 2 emails are same")))
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
	if err := render.Bind(r, Argument);err != nil{
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("invalid format")))
		return
	}
	if !repo.IsEmailValid(Argument.Requestor)||!repo.IsEmailValid(Argument.Target){
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("email is wrong")))
		return
	}
	if Argument.Requestor == Argument.Target {
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("error cause 2 emails are same")))
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
	if err := render.Bind(r, Argument);err != nil{
		render.Render(w,r,exception.ServerErrorRenderer(errors.New("invalid format")))
		return
	}
	if !repo.IsEmailValid(Argument.Sender){
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


