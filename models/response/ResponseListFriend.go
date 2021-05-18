package models

import "net/http"

type ResponseListFriend struct {
	Success bool     `json:"success"`
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}

func (*ResponseListFriend) Render(w http.ResponseWriter, r *http.Request) error{
	return nil
}