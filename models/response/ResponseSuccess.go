package models

import "net/http"

type ResponseSuccess struct {
	Success bool `json:"success"`
}

func (*ResponseSuccess) Render(w http.ResponseWriter, r *http.Request) error{
	return nil
}