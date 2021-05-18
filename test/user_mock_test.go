package test

import (
	contr "Friend_management/controller"
	"Friend_management/db"
	"Friend_management/models"
	repo "Friend_management/repository"
	ser "Friend_management/services"
	"bytes"
	"encoding/json"

	// "fmt"

	// "errors"
	"net/http"
	"net/http/httptest"
	"strings"

	// "bytes"
	// "database/sql"

	"errors"
	"testing"

	"io/ioutil"
	// gomock "github.com/golang/mock/gomock"
	// "github.com/sgreben/testing-with-gomock/mocks"
	// "github.com/sgreben/testing-with-gomock/user"
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
	//rs := args.Get(0)
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
func TestFindAll(t *testing.T) {
	mockRepo := new(MockRepository)
	u1 := models.User{Email: "len"}
	lst := make([]models.User, 0)
	lst = append(lst, u1)
	mockRepo.On("GetAllUsers").Return(&models.UserList{lst}, nil)
	testing := ser.NewUserService(mockRepo)
	rs, _ := testing.FindAllUser(db.Database{})
	mockRepo.AssertExpectations(t)
	assert.Equal(t, "len", rs.Users[0].Email)
}
func TestGetListUser(t *testing.T) {
	//given
	lst := &models.UserList{}
	u1 := &models.User{Email: "len1"}
	u2 := &models.User{Email: "len2"}
	u3 := &models.User{Email: "len3"}
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
			expectedSuccessBody: `{"users":[{"email":"len1"},{"email":"len2"},{"email":"len3"}]}`,
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
				y repo.UserRepoInter = mockRepo
				x ser.UserService    = ser.NewUserService(y)
			)
			req, err := http.NewRequest("GET", "/users", nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(contr.NewUserControl(x).GetAllUsers)
			handler.ServeHTTP(w, req)
			var Body models.UserList
			// var actualResult string
			if tc.scenario == "Success" {
				err = json.Unmarshal([]byte(w.Body.String()), &Body)
				if err != nil {
					t.Errorf("something wrong")
				}
				// body, _ := ioutil.ReadAll(w.Result().Body)
				// actualResult = string(body)
			}
			if tc.scenario == "Success" {
				assert.Equal(t, 200, w.Result().StatusCode)
				text, _ := w.Body.ReadString('\n')
				text = strings.Replace(text, "\n", "", -1)
				assert.Equal(t, tc.expectedSuccessBody, text)
			} else {
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, w.Body.String())
			}
		})
	}
}
func TestCreateNewUserController(t *testing.T) {
	CreateConnection()
	//given
	testCase := []struct {
		scenario          string
		inputRequest      *models.User
		expectedErrorBody string
	}{
		{
			scenario:          "Success",
			inputRequest:      &models.User{Email: "hcl2@gmail.com"},
			expectedErrorBody: "",
		},
		{
			scenario:          "Failure",
			inputRequest:      &models.User{Email: "hcl2@gmail.com"},
			expectedErrorBody: "error",
		},
		{
			scenario: 			"Wrong format email",
			inputRequest: 		&models.User{Email: "hcl"},
			expectedErrorBody: `"email was wrong"`,		
		},
		{
			scenario:          "Empty request body",
			expectedErrorBody: `"render: unable to automatically decode the request content type"`,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.scenario, func(t *testing.T) {
			mockRepo := new(MockRepository)
			if tc.inputRequest != nil{
				if tc.scenario != "Failure" {
					mockRepo.On("AddUser").Return(nil)
				}else{
					mockRepo.On("AddUser").Return(errors.New("Any error"))
				}
			}
			var (
				re       repo.UserRepoInter = mockRepo
				services ser.UserService    = ser.NewUserService(re)
			)
			var w *httptest.ResponseRecorder
			if tc.inputRequest != nil {
				value := map[string]string{"email": tc.inputRequest.Email}
				jsonValue, _ := json.Marshal(value)
				req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewUserControl(services).CreateUser)
				handler.ServeHTTP(w, req)
			} else {
				req, err := http.NewRequest("POST", "/users", nil)

				if err != nil {
					t.Fatal(err)
				}
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewUserControl(services).CreateUser)
				handler.ServeHTTP(w, req)
			}
			if tc.scenario == "Success" {
				assert.Equal(t, 200, w.Result().StatusCode)
			} else if tc.scenario == "Failure" {
				assert.Equal(t, 500, w.Result().StatusCode)
			} else if tc.scenario=="Wrong format email"{
				assert.Equal(t, tc.expectedErrorBody, w.Body.String())
			}else {
				body, _ := ioutil.ReadAll(w.Result().Body)
				assert.Equal(t, tc.expectedErrorBody, string(body))
			}
		})
	}
}
func TestDeleteUserController(t *testing.T) {
	CreateConnection()
	//given
	testCase := []struct {
		scenario          string
		inputRequest      *models.User
		expectedErrorBody string
	}{
		{
			scenario:          "Success",
			inputRequest:      &models.User{Email: "hcl2@gmail.com"},
			expectedErrorBody: "",
		},
		{
			scenario:          "Failure",
			inputRequest:      &models.User{Email: "hcl2@gmail.com"},
			expectedErrorBody: `"error"`,
		},
		{
			scenario:          "Wrong format email",
			inputRequest:      &models.User{Email: "hcl2"},
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
			if tc.inputRequest !=nil{
				if tc.scenario != "Failure"{
					mockRepo.On("DeleteUser").Return(nil)
				}else{
					mockRepo.On("DeleteUser").Return(errors.New("error"))
				}
			}

			var (
				re       repo.UserRepoInter = mockRepo
				services ser.UserService    = ser.NewUserService(re)
			)
			var w *httptest.ResponseRecorder
			if tc.inputRequest != nil {		
				req, err := http.NewRequest("DELETE", "/users/delete", nil)
				if err != nil {
					t.Fatal(err)
				}
				q := req.URL.Query()
				q.Add("id", tc.inputRequest.Email)
				req.URL.RawQuery = q.Encode()
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewUserControl(services).DeleteUser)
				handler.ServeHTTP(w, req)
			}else{
				req, err := http.NewRequest("DELETE", "/users/delete", nil)
				if err != nil {
					t.Fatal(err)
				}
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewUserControl(services).DeleteUser)
				handler.ServeHTTP(w, req)
			}
			if tc.scenario == "Success" {
				assert.Equal(t, 200, w.Result().StatusCode)
			} else if tc.scenario == "Wrong format email" {
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, w.Body.String())
			}else if tc.scenario =="Failure"{
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, w.Body.String())
			}else{
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, w.Body.String())
			}
		})
	}
}
func TestGetUserController(t *testing.T) {
	CreateConnection()
	//given
	user := models.User{Email: "len"}
	testCase := []struct {
		scenario          string
		inputRequest      *models.User
		expectedErrorBody string
	}{
		{
			scenario:          "Success",
			inputRequest:      &models.User{Email: "hcl2@gmail.com"},
			expectedErrorBody: "",
		},
		{
			scenario:          "Failure",
			inputRequest:      &models.User{Email: "hcl2@gmail.com"},
			expectedErrorBody: `"error"`,
		},
		{
			scenario:          "Wrong format email",
			inputRequest:      &models.User{Email: "hcl2"},
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
			if tc.inputRequest !=nil{
				if tc.scenario != "Failure"{
					mockRepo.On("GetUserByEmail").Return(user,nil)
				}else{
					mockRepo.On("GetUserByEmail").Return(user,errors.New("error"))
				}
			}
			var (
				re       repo.UserRepoInter = mockRepo
				services ser.UserService    = ser.NewUserService(re)
			)
			var w *httptest.ResponseRecorder
			if tc.inputRequest != nil {
				req, err := http.NewRequest("GET", "/users/find", nil)
				if err != nil {
					t.Fatal(err)
				}
				q := req.URL.Query()
				q.Add("id", tc.inputRequest.Email)
				req.URL.RawQuery = q.Encode()
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewUserControl(services).GetUser)
				handler.ServeHTTP(w, req)
			}else{
				req, err := http.NewRequest("GET", "/users/find", nil)
				if err != nil {
					t.Fatal(err)
				}
				w = httptest.NewRecorder()
				handler := http.HandlerFunc(contr.NewUserControl(services).GetUser)
				handler.ServeHTTP(w, req)
			}
			if tc.scenario == "Success" {
				assert.Equal(t, 200, w.Result().StatusCode)
			} else if tc.scenario == "Wrong format email" {
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, w.Body.String())
			}else if tc.scenario =="Failure"{
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, w.Body.String())
			}else{
				assert.Equal(t, 500, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, w.Body.String())
			}
		})
	}
}
//-------------------------------------------------------------------------------------