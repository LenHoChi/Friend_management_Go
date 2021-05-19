package test

import (
	contr "Friend_management/controller"
	"Friend_management/db"
	"Friend_management/models"
	r_Request "Friend_management/models/request"
	r_Response "Friend_management/models/response"
	repo "Friend_management/repository"
	ser "Friend_management/services"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)
type MockRepositoryRela struct {
	mock.Mock
}
func (mock *MockRepositoryRela)GetAllRelationship(database db.Database)(*models.RelationshipList, error){
	args := mock.Called()
	rs := args.Get(0)
	return rs.(*models.RelationshipList), args.Error(1)
}
func (mock *MockRepositoryRela)FindRelationshipByKey(database db.Database, userEmail string, friendEmail string) (models.Relationship, error){
	args := mock.Called()
	rs:= args.Get(0)
	return rs.(models.Relationship), args.Error(1)
}
func (mock *MockRepositoryRela)AddRelationship(database db.Database, userEmail string, friendEmail string) (*r_Response.ResponseSuccess, error){
	args := mock.Called()
	rs:= args.Get(0)
	return rs.(*r_Response.ResponseSuccess), args.Error(1)
}
func (mock *MockRepositoryRela)FindListFriend(database db.Database, email string) (*r_Response.ResponseListFriend, error){
	args := mock.Called()
	rs:= args.Get(0)
	return rs.(*r_Response.ResponseListFriend), args.Error(1)
}
func (mock *MockRepositoryRela)FindCommonListFriend(database db.Database, lstEmail []string) (*r_Response.ResponseListFriend, error){
	args := mock.Called()
	rs:= args.Get(0)
	return rs.(*r_Response.ResponseListFriend), args.Error(1)
}
func (mock *MockRepositoryRela)BeSubcribe(database db.Database, requestor string, target string) (*r_Response.ResponseSuccess, error){
	args := mock.Called()
	rs:= args.Get(0)
	return rs.(*r_Response.ResponseSuccess), args.Error(1)
}
func (mock *MockRepositoryRela)ToBlock(database db.Database, requestor string, target string) (*r_Response.ResponseSuccess, error){
	args := mock.Called()
	rs:= args.Get(0)
	return rs.(*r_Response.ResponseSuccess), args.Error(1)
}
func (mock *MockRepositoryRela)RetrieveUpdate(database db.Database, sender string, target string) (*r_Response.ResponseRetrieve, error){
	args := mock.Called()
	rs:= args.Get(0)
	return rs.(*r_Response.ResponseRetrieve), args.Error(1)
}
func TestGetListRelationship(t *testing.T){
	//given
	email1 := &models.User{Email: "hcl@gmail.com"}
	email2 := &models.User{Email: "hcl@gmail.com"}
	email3 := &models.User{Email: "hcl@gmail.com"}
	lst := &models.RelationshipList{}
	l1 := &models.Relationship{UserEmail: email1.Email, FriendEmail: email2.Email, AreFriend: true, IsSubcriber:  false, IsBlock:  false}
	l2 := &models.Relationship{UserEmail: email1.Email, FriendEmail: email3.Email, AreFriend: true, IsSubcriber:  false, IsBlock:  false}
	lst.Relationships = append(lst.Relationships, *l1, *l2)
	testCases := []struct{
		scenario string
		mockResponse *models.RelationshipList
		mockError error
		expectedErrorBody string
		expectedSuccessBody string
	}{
		{
			scenario: "Success",
			mockResponse: lst,
			expectedSuccessBody: `{"relationships":[{"UserEmail":"hcl@gmail.com","FriendEmail":"hcl@gmail.com","AreFriend":true,"IsSubcriber":false,"IsBlock":false},{"UserEmail":"hcl@gmail.com","FriendEmail":"hcl@gmail.com","AreFriend":true,"IsSubcriber":false,"IsBlock":false}]}`,
		},
		{
			scenario: "Failure",
			mockResponse:  lst,
			mockError:  errors.New("errors"),
			expectedErrorBody: `errors`,
		},
	}
	for _,tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			mockRepo := new(MockRepositoryRela)
			mockRepo.On("GetAllRelationship").Return(tc.mockResponse, tc.mockError)
			var (
				repoRelationship repo.RelationshipInter = mockRepo
				serviceRelationship ser.RepositoryService = ser.NewRelationshipService(repoRelationship)
			)
			req, err := http.NewRequest("GET","/relationship",nil)
			if err !=nil{
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(contr.NewRelationshipControl(serviceRelationship).GetAllRelationships)
			handler.ServeHTTP(w, req)
			var Body models.RelationshipList
			var Body2 r_Response.ResponseRenderError
			json.Unmarshal(w.Body.Bytes(), &Body)
			json.Unmarshal(w.Body.Bytes(),&Body2)
			if tc.scenario == "Success"{
				assert.Equal(t, 200, w.Result().StatusCode)
				text := w.Body.String()
				text = strings.Replace(text, "\n", "", -1)
				assert.Equal(t, tc.expectedSuccessBody, text)
				assert.Equal(t, "hcl@gmail.com", Body.Relationships[0].UserEmail)
			}else{
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody,Body2.Message)
				assert.Equal(t, "Internal server error",Body2.StatusText)
			}
		})
	}
}
func TestMakeFriendController(t *testing.T){
	db.Initialize()
	//given
	testCases := []struct{
		scenario string
		inputRequest *r_Request.RequestFriendLists
		expectedErrorBody string
		mockResponse *r_Response.ResponseSuccess
	}{
		{
			scenario: "Success",
			inputRequest: &r_Request.RequestFriendLists{
				RequestFriendLists: []string{
					"hcl@gmail.com",
					"hcl2@gmail.com",
				},
			},
			mockResponse: &r_Response.ResponseSuccess{Success: true},
			expectedErrorBody: "",
		},
		{
			scenario: "Failure",
			inputRequest: &r_Request.RequestFriendLists{
				RequestFriendLists: []string{
					"hcl@gmail.com",
					"hcl2@gmail.com",
				},
			},
			expectedErrorBody: "error",
		},
		{
			scenario: "Wrong format email",
			inputRequest: &r_Request.RequestFriendLists{
				RequestFriendLists: []string{
					"hcl",
					"hcl2",
				},
			},
			expectedErrorBody: `email is wrong`,
		},
	}
	for _, tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			mockRepo := new(MockRepositoryRela)
			if tc.inputRequest != nil{
				if tc.scenario!= "Failure"{
					mockRepo.On("AddRelationship").Return(tc.mockResponse,nil)
				}else{
					mockRepo.On("AddRelationship").Return(tc.mockResponse,errors.New("error"))
				}
			}
			var (
				repoRelationship repo.RelationshipInter = mockRepo
				serviceRelationship ser.RepositoryService = ser.NewRelationshipService(repoRelationship)
			)
			var Body2 r_Response.ResponseRenderError
			var w *httptest.ResponseRecorder
			if tc.inputRequest != nil {
				values := map[string][]string{"friends": tc.inputRequest.RequestFriendLists}
				jsonValue,_ := json.Marshal(values)
				req, err := http.NewRequest("POST","/relationship/make",bytes.NewBuffer(jsonValue))
				if err != nil{
					t.Fatal(err)
				}
				req.Header.Set("Content-Type","application/json")
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewRelationshipControl(serviceRelationship).MakeFriend)
				handler.ServeHTTP(w, req)
				json.Unmarshal(w.Body.Bytes(),&Body2)
			}
			if tc.scenario == "Success"{
				assert.Equal(t, 200, w.Result().StatusCode)
			}else if tc.scenario == "Failure"{
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, Body2.Message)
			}else if tc.scenario =="Wrong format email"{
				assert.Equal(t, tc.expectedErrorBody, Body2.Message)
			}
		})
	}
}
func TestFindListFriendController(t *testing.T){
	//given
	testCases := []struct{
		scenario string
		inputRequest *r_Request.RequestEmail
		expectedErrorBody string
		mockResponse *r_Response.ResponseListFriend
	}{
		{
			scenario: "Success",
			inputRequest: &r_Request.RequestEmail{
				Email: "hcl@gmail.com",
			},
			mockResponse: &r_Response.ResponseListFriend{
				Success: true,
				Friends: []string{
					"hcl1@gmail.com",
					"hcl2@gmail.com",
				},
				Count: 2,
			},
			expectedErrorBody: "",
		},
		{
			scenario: "Failure",
			inputRequest: &r_Request.RequestEmail{
				Email: "hcl@gmail.com",
			},
			expectedErrorBody: "error",
		},
		{
			scenario: "Wrong format email",
			inputRequest: &r_Request.RequestEmail{
				Email: "hcl",
			},
			expectedErrorBody: `email is wrong`,
		},
	}
	for _, tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			mockRepo := new(MockRepositoryRela)
			if tc.inputRequest != nil{
				if tc.scenario!= "Failure"{
					mockRepo.On("FindListFriend").Return(tc.mockResponse,nil)
				}else{
					mockRepo.On("FindListFriend").Return(tc.mockResponse,errors.New("error"))
				}
			}
			var (
				repoRelationship repo.RelationshipInter = mockRepo
				serviceRelationship ser.RepositoryService = ser.NewRelationshipService(repoRelationship)
			)
			var Body2 r_Response.ResponseRenderError
			var w *httptest.ResponseRecorder
			if tc.inputRequest != nil {
				value := map[string]string{"email": tc.inputRequest.Email}
				jsonValue,_ := json.Marshal(value)
				req, err := http.NewRequest("POST","/relationship/list",bytes.NewBuffer(jsonValue))
				if err != nil{
					t.Fatal(err)
				}
				req.Header.Set("Content-Type","application/json")
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewRelationshipControl(serviceRelationship).FindListFriend)
				handler.ServeHTTP(w, req)
				json.Unmarshal(w.Body.Bytes(),&Body2)
			}
			if tc.scenario == "Success"{
				assert.Equal(t, 200, w.Result().StatusCode)
			}else if tc.scenario == "Failure"{
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, Body2.Message)
			}else if tc.scenario =="Wrong format email"{
				assert.Equal(t, tc.expectedErrorBody, Body2.Message)
			}
		})
	}
}
func TestFindCommonListFriendController(t *testing.T){
	//given
	testCases := []struct{
		scenario string
		inputRequest *r_Request.RequestFriendLists
		expectedErrorBody string
		mockResponse *r_Response.ResponseListFriend
	}{
		{
			scenario: "Success",
			inputRequest: &r_Request.RequestFriendLists{
				RequestFriendLists: []string{
					"hcl@gmail.com",
					"hcl2@gmail.com",
				},
			},
			mockResponse: &r_Response.ResponseListFriend{
				Success: true,
				Friends: []string{
					"hcl5@gmail.com",
					"hcl6@gmail.com",
				},
				Count: 2,
			},
			expectedErrorBody: "",
		},
		{
			scenario: "Failure",
			inputRequest: &r_Request.RequestFriendLists{
				RequestFriendLists: []string{
					"hcl@gmail.com",
					"hcl2@gmail.com",
				},
			},
			expectedErrorBody: "error",
		},
		{
			scenario: "Wrong format email",
			inputRequest: &r_Request.RequestFriendLists{
				RequestFriendLists: []string{
					"hcl",
					"hcl2",
				},
			},
			expectedErrorBody: `email is wrong`,
		},
	}
	for _, tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			mockRepo := new(MockRepositoryRela)
			if tc.inputRequest != nil{
				if tc.scenario!= "Failure"{
					mockRepo.On("FindCommonListFriend").Return(tc.mockResponse,nil)
				}else{
					mockRepo.On("FindCommonListFriend").Return(tc.mockResponse,errors.New("error"))
				}
			}
			var (
				repoRelationship repo.RelationshipInter = mockRepo
				serviceRelationship ser.RepositoryService = ser.NewRelationshipService(repoRelationship)
			)
			var Body2 r_Response.ResponseRenderError
			var w *httptest.ResponseRecorder
			if tc.inputRequest != nil {
				jsonValue,_ := json.Marshal(tc.inputRequest)
				req, err := http.NewRequest("POST","/relationship/common",bytes.NewBuffer(jsonValue))
				if err != nil{
					t.Fatal(err)
				}
				req.Header.Set("Content-Type","application/json")
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewRelationshipControl(serviceRelationship).FindCommonListFriend)
				handler.ServeHTTP(w, req)
				json.Unmarshal(w.Body.Bytes(),&Body2)
			}
			if tc.scenario == "Success"{
				assert.Equal(t, 200, w.Result().StatusCode)
			}else if tc.scenario == "Failure"{
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, Body2.Message)
			}else if tc.scenario =="Wrong format email"{
				assert.Equal(t, tc.expectedErrorBody, Body2.Message)
			}
		})
	}
}
func TestBeSubcriberController(t *testing.T){
	db.Initialize()
	//given
	testCases := []struct{
		scenario string
		inputRequest *r_Request.RequestUpdate
		expectedErrorBody string
		mockResponse *r_Response.ResponseSuccess
	}{
		{
			scenario: "Success",
			inputRequest: &r_Request.RequestUpdate{Requestor: "hcl@gmail.com", Target: "hcl1@gmail.com"},
			mockResponse: &r_Response.ResponseSuccess{Success: true},
			expectedErrorBody: "",
		},
		{
			scenario: "Failure",
			inputRequest: &r_Request.RequestUpdate{Requestor: "hcl@gmail.com", Target: "hcl1@gmail.com"},
			expectedErrorBody: "error",
		},
		{
			scenario: "Wrong format email",
			inputRequest: &r_Request.RequestUpdate{Requestor: "hcl", Target: "hcl1@gmail.com"},
			expectedErrorBody: `email is wrong`,
		},
	}
	for _, tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			mockRepo := new(MockRepositoryRela)
			if tc.inputRequest != nil{
				if tc.scenario!= "Failure"{
					mockRepo.On("BeSubcribe").Return(tc.mockResponse,nil)
				}else{
					mockRepo.On("BeSubcribe").Return(tc.mockResponse,errors.New("error"))
				}
			}
			var (
				repoRelationship repo.RelationshipInter = mockRepo
				serviceRelationship ser.RepositoryService = ser.NewRelationshipService(repoRelationship)
			)
			var Body2 r_Response.ResponseRenderError
			var w *httptest.ResponseRecorder
			if tc.inputRequest != nil {
				jsonValue,_ := json.Marshal(tc.inputRequest)
				req, err := http.NewRequest("POST","/relationship/update",bytes.NewBuffer(jsonValue))
				if err != nil{
					t.Fatal(err)
				}
				req.Header.Set("Content-Type","application/json")
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewRelationshipControl(serviceRelationship).BeSubcriber)
				handler.ServeHTTP(w, req)
				json.Unmarshal(w.Body.Bytes(),&Body2)
			}
			if tc.scenario == "Success"{
				assert.Equal(t, 200, w.Result().StatusCode)
			}else if tc.scenario == "Failure"{
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, Body2.Message)
			}else if tc.scenario =="Wrong format email"{
				assert.Equal(t, tc.expectedErrorBody, Body2.Message)
			}
		})
	}
}
func TestToBlockController(t *testing.T){
	db.Initialize()
	//given
	testCases := []struct{
		scenario string
		inputRequest *r_Request.RequestUpdate
		expectedErrorBody string
		mockResponse *r_Response.ResponseSuccess
	}{
		{
			scenario: "Success",
			inputRequest: &r_Request.RequestUpdate{Requestor: "hcl@gmail.com", Target: "hcl1@gmail.com"},
			mockResponse: &r_Response.ResponseSuccess{Success: true},
			expectedErrorBody: "",
		},
		{
			scenario: "Failure",
			inputRequest: &r_Request.RequestUpdate{Requestor: "hcl@gmail.com", Target: "hcl1@gmail.com"},
			expectedErrorBody: "error",
		},
		{
			scenario: "Wrong format email",
			inputRequest: &r_Request.RequestUpdate{Requestor: "hcl", Target: "hcl1@gmail.com"},
			expectedErrorBody: `email is wrong`,
		},
	}
	for _, tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			mockRepo := new(MockRepositoryRela)
			if tc.inputRequest != nil{
				if tc.scenario!= "Failure"{
					mockRepo.On("ToBlock").Return(tc.mockResponse,nil)
				}else{
					mockRepo.On("ToBlock").Return(tc.mockResponse,errors.New("error"))
				}
			}
			var (
				repoRelationship repo.RelationshipInter = mockRepo
				serviceRelationship ser.RepositoryService = ser.NewRelationshipService(repoRelationship)
			)
			var Body2 r_Response.ResponseRenderError
			var w *httptest.ResponseRecorder
			if tc.inputRequest != nil {
				jsonValue,_ := json.Marshal(tc.inputRequest)
				req, err := http.NewRequest("POST","/relationship/block",bytes.NewBuffer(jsonValue))
				if err != nil{
					t.Fatal(err)
				}
				req.Header.Set("Content-Type","application/json")
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewRelationshipControl(serviceRelationship).ToBLock)
				handler.ServeHTTP(w, req)
				json.Unmarshal(w.Body.Bytes(),&Body2)
			}
			if tc.scenario == "Success"{
				assert.Equal(t, 200, w.Result().StatusCode)
			}else if tc.scenario == "Failure"{
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, Body2.Message)
			}else if tc.scenario =="Wrong format email"{
				assert.Equal(t, tc.expectedErrorBody, Body2.Message)
			}
		})
	}
}
func TestRetrieveUpdateController(t *testing.T){
	//given
	testCases := []struct{
		scenario string
		inputRequest *r_Request.RetrieveUpdate
		expectedErrorBody string
		mockResponse *r_Response.ResponseRetrieve
	}{
		{
			scenario: "Success",
			inputRequest: &r_Request.RetrieveUpdate{Sender: "hcl@gmail.com", Tartget: "hcl1@gmail.com"},
			mockResponse: &r_Response.ResponseRetrieve{Success: true, Recipients: []string{"hcl2@gmail.com, hcl1@gmail.com"}},
			expectedErrorBody: "",
		},
		{
			scenario: "Failure",
			inputRequest: &r_Request.RetrieveUpdate{Sender: "hcl@gmail.com", Tartget: "hcl1@gmail.com"},
			expectedErrorBody: "error",
		},
		{
			scenario: "Wrong format email",
			inputRequest: &r_Request.RetrieveUpdate{Sender: "hcl", Tartget: "hcl1@gmail.com"},
			expectedErrorBody: `email is wrong`,
		},
	}
	for _, tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			mockRepo := new(MockRepositoryRela)
			if tc.inputRequest != nil{
				if tc.scenario!= "Failure"{
					mockRepo.On("RetrieveUpdate").Return(tc.mockResponse,nil)
				}else{
					mockRepo.On("RetrieveUpdate").Return(tc.mockResponse,errors.New("error"))
				}
			}
			var (
				repoRelationship repo.RelationshipInter = mockRepo
				serviceRelationship ser.RepositoryService = ser.NewRelationshipService(repoRelationship)
			)
			var Body2 r_Response.ResponseRenderError
			var w *httptest.ResponseRecorder
			if tc.inputRequest != nil {
				jsonValue,_ := json.Marshal(tc.inputRequest)
				req, err := http.NewRequest("POST","/relationship/retrieve",bytes.NewBuffer(jsonValue))
				if err != nil{
					t.Fatal(err)
				}
				req.Header.Set("Content-Type","application/json")
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewRelationshipControl(serviceRelationship).RetrieveUpdate)
				handler.ServeHTTP(w, req)
				json.Unmarshal(w.Body.Bytes(),&Body2)
			}
			if tc.scenario == "Success"{
				assert.Equal(t, 200, w.Result().StatusCode)
			}else if tc.scenario == "Failure"{
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, Body2.Message)
			}else if tc.scenario =="Wrong format email"{
				assert.Equal(t, tc.expectedErrorBody, Body2.Message)
			}
		})
	}
}