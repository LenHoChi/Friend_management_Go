package models

import (
	"fmt"
	"net/http"
)

type User struct {
	Email string `json:"email"`
}

type UserList struct {
	Users []User `json:"users"`
}

// func (*UserList) Render(w http.ResponseWriter, r *http.Request) error {
// 	return nil
// }
func (UserList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u *User) Bind(r *http.Request) error {
	if u.Email == "" {
		return fmt.Errorf("email is a required field")
	}
	return nil
}
