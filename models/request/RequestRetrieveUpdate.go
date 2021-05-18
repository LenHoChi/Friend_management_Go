package models
import ("net/http")
type RetrieveUpdate struct{
	Sender string `json:"sender"`
	Tartget string `json:"target"`
}
type RetrieveUpdateList struct{
	RetrieveUpdateList []RetrieveUpdate `json:"retrieveupdatelist"`
}
func (*RetrieveUpdate) Bind(r *http.Request) error{
	return nil
}
func (*RetrieveUpdateList) Render(w http.ResponseWriter, r *http.Request)error{
	return nil
}