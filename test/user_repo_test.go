package test

import (
	contr "Friend_management/controller"
	"Friend_management/models"
	"Friend_management/repository"
	"errors"
	"testing"
	"net/http"
	"context"
	"log"
	"net"
	"Friend_management/db"
	"github.com/stretchr/testify/assert"
	"os"
	"time"
)
func TestGetUser(t *testing.T){
	CreateConnection()
	user := &models.User{Email: "hcl@gmail.com"}
	repository.NewRepo().AddUser(contr.DBInstance, user)
	testCases := []struct{
		scenario	string
		mockInput	models.User
		expectedError error
	}{
		{
			scenario: "Success",
			mockInput: *user,
			expectedError: nil,
		},
	}
	for _, tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			_,actualRs := repository.NewRepo().GetUserByEmail(contr.DBInstance,tc.mockInput.Email)
			assert.Equal(t, tc.expectedError, actualRs)
		})
	}
}
func TestDeleteUser(t *testing.T){
	CreateConnection()
	//add user
	user := &models.User{Email: "hcl@gmail.com"}
	user3 := &models.User{Email: "hcl2@gmail.com"}
	repository.NewRepo().AddUser(contr.DBInstance,user)

	testCases := []struct {
		scenario	string
		mockInput	 string
		expectedError	 error
	}{
		{
			scenario:	"Success",
			mockInput: user.Email,
			expectedError: nil,
		},
		{
			scenario: "User not exists",
			mockInput: user3.Email,
			expectedError: errors.New("this user not exists"),
		},
	}

	for _, tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			actualRs := repository.NewRepo().DeleteUser(contr.DBInstance,tc.mockInput)
			assert.Equal(t, tc.expectedError, actualRs)
		})
	}
}
func TestCreateNewUser (t *testing.T){
	const numsUser int = 1
	lstUsers := &models.UserList{}
	for i := 0;i<numsUser;i++{
		user := &models.User{Email: "hcl@gmail.com"}
		lstUsers.Users = append(lstUsers.Users, *user)
	}
	CreateConnection()
	testCases := []struct{
		scenario	string
		mockInput models.User
		expectedError error
	}{
		{
			scenario: "Success",
			mockInput: lstUsers.Users[0],
			expectedError: nil,
		},
		{
			scenario: 	"User Exists",
			mockInput: 	lstUsers.Users[0],
			expectedError: errors.New("this email exists already"),
		},
	}

	for _, tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			actualRs := repository.NewRepo().AddUser(contr.DBInstance,&tc.mockInput)
			assert.Equal(t, tc.expectedError, actualRs)
		})
	}
}
func TestGetAllListUsers(t *testing.T){
	CreateConnection()

	result, err := repository.NewRepo().GetAllUsers(contr.DBInstance)
	assert.Nil(t, err)
	assert.NotNil(t, result)
}
func CreateConnection() {
	addr := ":8080"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error occurred: %s", err.Error())
	}
	database, err := db.Initialize2()
	// out : Conn &*sql.DB
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}
	// defer database.Conn.Close()
	//input : Conn &*sql.Tx
	// db2,_ := database.Conn.Begin()
	// httpHandler := contr.NewHandler2(db2)
	httpHandler := contr.NewHandler(database)
	server := &http.Server{
		Handler: httpHandler,
	}
	go func() {
		server.Serve(listener)
	}()
	defer Stop(server)
}
func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Could not shut down server correctly: %v\n", err)
		os.Exit(1)
	}
}