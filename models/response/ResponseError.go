package models

import (
	"encoding/json"
	"net/http"
)

type ResponseError struct{
	Code int `json:"-"`
	Description string `json:"description"`
}
type ResponseRenderError struct{
	StatusText string `json:"status_text"`
	Message string `json:"message"`
}
// func (e *ResponseError) Error() string{
// 	return fmt.Sprintf("%s", e.Description)
// }

func ResponseWithJSON(response http.ResponseWriter, statusCode int, data interface{}){
	result, _ := json.Marshal(data)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	response.Write(result)
}
func ResponseWithError(response http.ResponseWriter, statusCode int, msg string){
	ResponseWithJSON(response, statusCode, map[string]string{
		"error":msg,
	})
}