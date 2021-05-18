package models

import (
	"net/http")
type RequestUpdate struct{
	Requestor string `json:"requestor"`
	Target string `json:"target"`
}
func (u *RequestUpdate) Bind(r *http.Request) error{
	return nil
}