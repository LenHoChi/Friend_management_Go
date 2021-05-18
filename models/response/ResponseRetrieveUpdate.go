package models

import "net/http"

type ResponseRetrieve struct{
	Success bool `json:"success"`
	Recipients []string `json:"recipients"`
}
func (*ResponseRetrieve) Render(w http.ResponseWriter, r *http.Request) error{
	return nil
}
