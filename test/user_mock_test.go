package test

import (
	contr "Friend_management/controller"
	"Friend_management/db"
	"Friend_management/models"
	repo "Friend_management/repository"
	ser "Friend_management/services"
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"strings"

	"errors"
	"testing"

	// "io/ioutil"
	// "io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mock *MockRepository) GetAllUsers(database db.Database) (*models.UserList, error) {
	args := mock.Called()
	rs := args.Get(0)
	return rs.(*models.UserList), args.Error(1)
}
func (mock *MockRepository) AddUser(database db.Database, user *models.User) error {
	args := mock.Called()
	return args.Error(0)
}
func (mock *MockRepository) GetUserByEmail(database db.Database, email string) (models.User, error) {
	args := mock.Called()
	rs := args.Get(0)
	return rs.(models.User), args.Error(1)
}
func (mock *MockRepository) DeleteUser(database db.Database, email string) error {
	args := mock.Called()
	return args.Error(0)
}
func TestGetListUser(t *testing.T) {
	//given
	lst := &models.UserList{}
	u1 := &models.User{Email: "hcl@gmail.com"}
	u2 := &models.User{Email: "hcl2@gmail.com"}
	u3 := &models.User{Email: "hcl3@gmail.com"}
	lst.Users = append(lst.Users, *u1, *u2, *u3)
	testCases := []struct {
		scenario            string
		mockResponse        *models.UserList
		mockError           error
		expectedErrorBody   string
		expectedSuccessBody string
	}{
		{
			scenario:            "Success",
			mockResponse:        lst,
			expectedSuccessBody: `{"users":[{"email":"hcl@gmail.com"},{"email":"hcl2@gmail.com"},{"email":"hcl3@gmail.com"}]}`,
			//mockerorr = nil --> return = nil -- controller -->success
		},
		{
			scenario:          "Failuire",
			mockResponse:      lst,
			mockError:         errors.New("errors"),
			expectedErrorBody: `"errors"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockRepo.On("GetAllUsers").Return(tc.mockResponse, tc.mockError)

			var (
				repoUser     repo.UserRepoInter = mockRepo
				servicesUser ser.UserService    = ser.NewUserService(repoUser)
			)
			req, err := http.NewRequest("GET", "/users", nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(contr.NewUserControl(servicesUser).GetAllUsers)
			handler.ServeHTTP(w, req)
			var Body models.UserList
			if tc.scenario == "Success" {
				err = json.Unmarshal(w.Body.Bytes(), &Body)
				if err != nil {
					t.Errorf("something wrong")
				}
				assert.Equal(t, 200, w.Result().StatusCode)
				text, _ := w.Body.ReadString('\n')
				text = strings.Replace(text, "\n", "", -1)
				assert.Equal(t, tc.expectedSuccessBody, text)
				assert.Equal(t, u1.Email, Body.Users[0].Email)
			} else {
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, w.Body.String())
			}
		})
	}
}
func TestCreateNewUserController(t *testing.T) {
	db.Initialize()
	//given
	testCase := []struct {
		scenario          string
		inputRequest      *models.User
		expectedErrorBody string
	}{
		{
			scenario:          "Success",
			inputRequest:      &models.User{Email: "hcl@gmail.com"},
			expectedErrorBody: "",
		},
		{
			scenario:          "Failure",
			inputRequest:      &models.User{Email: "hcl@gmail.com"},
			expectedErrorBody: `"error"`,
		},
		{
			scenario:          "Wrong format email",
			inputRequest:      &models.User{Email: "hcl"},
			expectedErrorBody: `"email was wrong"`,
		},
		{
			scenario:          "Empty request body",
			inputRequest:      &models.User{Email: ""},
			expectedErrorBody: `"email is a required field"`,
		},
		{
			scenario:          "Empty request body",
			expectedErrorBody: `"render: unable to automatically decode the request content type"`,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.scenario, func(t *testing.T) {
			mockRepo := new(MockRepository)
			if tc.inputRequest != nil {
				if tc.scenario != "Failure" {
					mockRepo.On("AddUser").Return(nil)
				} else {
					mockRepo.On("AddUser").Return(errors.New("error"))
				}
			}
			var (
				repoUser     repo.UserRepoInter = mockRepo
				servicesUser ser.UserService    = ser.NewUserService(repoUser)
			)
			var w *httptest.ResponseRecorder
			if tc.inputRequest !=nil{
				value := map[string]string{"email": tc.inputRequest.Email}
				jsonValue, _ := json.Marshal(value)
				req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewUserControl(servicesUser).CreateUser)
				handler.ServeHTTP(w, req)
			}else{
				req, err := http.NewRequest("POST", "/users", nil)
				if err != nil {
					t.Fatal(err)
				}
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewUserControl(servicesUser).CreateUser)
				handler.ServeHTTP(w, req)
			}
			if tc.scenario == "Success" {
				assert.Equal(t, 200, w.Result().StatusCode)
			} else{
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, w.Body.String())
			}
		})
	}
}
func TestDeleteUserController(t *testing.T) {
	db.Initialize()
	//given
	testCase := []struct {
		scenario          string
		inputRequest      *models.User
		expectedErrorBody string
	}{
		{
			scenario:          "Success",
			inputRequest:      &models.User{Email: "hcl@gmail.com"},
			expectedErrorBody: "",
		},
		{
			scenario:          "Failure",
			inputRequest:      &models.User{Email: "hcl@gmail.com"},
			expectedErrorBody: `"error"`,
		},
		{
			scenario:          "Wrong format email",
			inputRequest:      &models.User{Email: "hcl"},
			expectedErrorBody: `"email was wrong"`,
		},
		{
			scenario:          "Empty request body",
			inputRequest:      &models.User{Email: ""},
			expectedErrorBody: `"email was wrong"`,
		},
		{
			scenario:          "Empty request body",
			expectedErrorBody: `"email was wrong"`,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.scenario, func(t *testing.T) {
			mockRepo := new(MockRepository)
			if tc.inputRequest != nil {
				if tc.scenario != "Failure" {
					mockRepo.On("DeleteUser").Return(nil)
				} else {
					mockRepo.On("DeleteUser").Return(errors.New("error"))
				}
			}
			var (
				repoUser     repo.UserRepoInter = mockRepo
				servicesUser ser.UserService    = ser.NewUserService(repoUser)
			)
			var w *httptest.ResponseRecorder
			if tc.inputRequest != nil{
				req, err := http.NewRequest("DELETE", "/users/delete", nil)
				if err != nil {
					t.Fatal(err)
				}
				q := req.URL.Query()
				q.Add("id", tc.inputRequest.Email)
				req.URL.RawQuery = q.Encode()
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewUserControl(servicesUser).DeleteUser)
				handler.ServeHTTP(w, req)
			}else{
				req, err := http.NewRequest("DELETE", "/users/delete", nil)
				if err != nil {
					t.Fatal(err)
				}
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewUserControl(servicesUser).DeleteUser)
				handler.ServeHTTP(w, req)
			}
			if tc.scenario == "Success" {
				assert.Equal(t, 200, w.Result().StatusCode)
			} else{
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, w.Body.String())
			}
		})
	}
}
func TestGetUserController(t *testing.T) {
	//given
	user := models.User{Email: "hcl@gmail.com"}
	testCase := []struct {
		scenario          string
		inputRequest      *models.User
		expectedErrorBody string
	}{
		{
			scenario:          "Success",
			inputRequest:      &models.User{Email: "hcl@gmail.com"},
			expectedErrorBody: "",
		},
		{
			scenario:          "Failure",
			inputRequest:      &models.User{Email: "hcl@gmail.com"},
			expectedErrorBody: `"error"`,
		},
		{
			scenario:          "Wrong format email",
			inputRequest:      &models.User{Email: "hcl"},
			expectedErrorBody: `"email was wrong"`,
		},
		{
			scenario:          "Empty request body",
			inputRequest:      &models.User{Email: ""},
			expectedErrorBody: `"email was wrong"`,
		},
		{
			scenario:          "Empty request body",
			expectedErrorBody: `"email was wrong"`,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.scenario, func(t *testing.T) {
			mockRepo := new(MockRepository)
			if tc.inputRequest != nil {
				if tc.scenario != "Failure" {
					mockRepo.On("GetUserByEmail").Return(user, nil)
				} else {
					mockRepo.On("GetUserByEmail").Return(user, errors.New("error"))
				}
			}
			var (
				repoUser     repo.UserRepoInter = mockRepo
				servicesUser ser.UserService    = ser.NewUserService(repoUser)
			)
			var w *httptest.ResponseRecorder
			var Body = &models.User{}
			if tc.inputRequest != nil{
				req, err := http.NewRequest("GET", "/users/find", nil)
				if err != nil {
					t.Fatal(err)
				}
				q := req.URL.Query()
				q.Add("id", tc.inputRequest.Email)
				req.URL.RawQuery = q.Encode()
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewUserControl(servicesUser).GetUser)
				handler.ServeHTTP(w, req)
				json.Unmarshal(w.Body.Bytes(),&Body)
			} else {
				req, err := http.NewRequest("GET", "/users/find", nil)
				if err != nil {
					t.Fatal(err)
				}
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewUserControl(servicesUser).GetUser)
				handler.ServeHTTP(w, req)
			}
			// json.NewDecoder(io.Reader(w.Body)).Decode(&Body)
			if tc.scenario == "Success" {
				assert.Equal(t, 200, w.Result().StatusCode)
				assert.Equal(t, tc.inputRequest.Email, Body.Email)
			} else{
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, w.Body.String())
			}
		})
	}
}
