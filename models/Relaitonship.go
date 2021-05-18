package models

import (
	"net/http"
	"fmt"
)

type Relationship struct {
	UserEmail string `json"user_email"`
	FriendEmail string `json"friend_email"`
	AreFriend bool `json"arefriends"`
	IsSubcriber bool `json"issubcriber"`
	IsBlock bool `json"isblock"`
} 

type RelationshipList struct {
	Relationships []Relationship `json:"relationships"`
}

func (*RelationshipList) Render (w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*Relationship) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u *Relationship) Bind(r *http.Request) error {
	if u.UserEmail == ""{
		return fmt.Errorf("email is a required field")
	}
	return nil
}