package models

import (
	"net/http" 
)

type RequestEmail struct{
	Email string `json:"email"`
}
type RequestFriendLists struct {
	RequestFriendLists []string `json:"friends"`
}
func (u *RequestEmail) Bind(r *http.Request) error{
	return nil
}
func (u *RequestFriendLists) Bind(r *http.Request) error{
	return nil
}
func (u *RequestFriendLists) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}