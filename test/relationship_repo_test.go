package test

import (
	"Friend_management/models"
	"Friend_management/repository"
	"errors"
	"testing"
	r_Request "Friend_management/models/request"
	"github.com/stretchr/testify/assert"
	"Friend_management/util"
)
func TestGetAllRelationship(t *testing.T){
	CreateConnection()

	result, err := repository.NewRepoRelationship().GetAllRelationship(util.DBInstance)
	assert.Nil(t, err)
	assert.NotNil(t, result)
}
func TestFindRelationshipByKey(t *testing.T){
	CreateConnection()
	
	rela := &models.Relationship{UserEmail: "hcl@gmail.com", FriendEmail: "pvq@gmail.com", AreFriend:  true, IsSubcriber:  false, IsBlock:  false}
	repository.NewRepo().AddUser(util.DBInstance,&models.User{Email: rela.UserEmail})
	repository.NewRepo().AddUser(util.DBInstance,&models.User{Email: rela.FriendEmail})
	repository.NewRepoRelationship().AddRelationship(util.DBInstance, rela.UserEmail, rela.FriendEmail)
	result, err := repository.NewRepoRelationship().FindRelationshipByKey(util.DBInstance, rela.UserEmail, rela.FriendEmail)
	assert.Nil(t, err)
	assert.NotNil(t, result)
}
func TestAddRelationship(t *testing.T){
	CreateConnection()
	lst := make([]string,0)
	lst = append(lst, "hcl@gmail.com", "pvq@gmail.com")
	lst3 := make([]string,0)
	lst3 = append(lst3, "ngv@gmail.com", "pvq@gmail.com")
	repository.NewRepo().AddUser(util.DBInstance,&models.User{Email: "hcl@gmail.com"})
	repository.NewRepo().AddUser(util.DBInstance,&models.User{Email: "pvq@gmail.com"})
	testCases := []struct{
		scenario string
		mockInput r_Request.RequestFriendLists
		expectedError error
	}{
		{
			scenario: "Success",
			mockInput: r_Request.RequestFriendLists{RequestFriendLists :lst},
			expectedError: nil,
		},
		{
			scenario: "Fail by not exists",
			mockInput: r_Request.RequestFriendLists{RequestFriendLists :lst3},
			expectedError: errors.New("user not exists"),
		},
		{
			scenario: "Fail by exists",
			mockInput: r_Request.RequestFriendLists{RequestFriendLists :lst},
			expectedError: errors.New("this relationship exists already"),
		},
	}
	for _, tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			_, actualRs := repository.NewRepoRelationship().AddRelationship(util.DBInstance, tc.mockInput.RequestFriendLists[0], tc.mockInput.RequestFriendLists[1])
			assert.Equal(t, tc.expectedError, actualRs)
		})
	}
}
func TestFindListFriend(t *testing.T){
	CreateConnection()
	user1 := &models.User{Email: "hcl@gmail.com"}
	user2 := &models.User{Email: "pvq@gmail.com"}
	user3 := &models.User{Email: "ntt@gmail.com"}
	user4 := &models.User{Email: "vvh@gmail.com"}
	repository.NewRepo().AddUser(util.DBInstance, user1)
	repository.NewRepo().AddUser(util.DBInstance, user2)
	repository.NewRepo().AddUser(util.DBInstance, user3)
	repository.NewRepoRelationship().AddRelationship(util.DBInstance, user1.Email, user2.Email)
	repository.NewRepoRelationship().AddRelationship(util.DBInstance, user1.Email, user3.Email)
	testCases :=[]struct{
		scenario string
		mockInput r_Request.RequestEmail
		count int
		Success bool
		expectedError error
	}{
		{
			scenario: "Success",
			mockInput: r_Request.RequestEmail{Email: user1.Email},
			count: 2,
			Success: true,
			expectedError: nil,
		},
		{
			scenario: "Fail by no user",
			mockInput: r_Request.RequestEmail{Email: user4.Email},
			count: 0,
			expectedError: errors.New("no users in table"),
		},
	}
	for _, tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			result, err := repository.NewRepoRelationship().FindListFriend(util.DBInstance, tc.mockInput.Email)
			if tc.scenario == "Success"{
				assert.Equal(t, tc.count, result.Count)
				assert.Equal(t, tc.expectedError, err)
				assert.True(t, tc.Success)
			}else{
				assert.Nil(t, result)
				assert.Equal(t, tc.expectedError, err)
			}
		})
	}
}
func TestFindCommonListFriend(t *testing.T){
	CreateConnection()
	user1 := &models.User{Email: "hcl@gmail.com"}
	user2 := &models.User{Email: "pvq@gmail.com"}
	user3 := &models.User{Email: "ntt@gmail.com"}
	repository.NewRepo().AddUser(util.DBInstance, user1)
	repository.NewRepo().AddUser(util.DBInstance, user2)
	repository.NewRepo().AddUser(util.DBInstance, user3)
	repository.NewRepoRelationship().AddRelationship(util.DBInstance, user1.Email, user2.Email)
	repository.NewRepoRelationship().AddRelationship(util.DBInstance, user3.Email, user2.Email)
	lst := make([]string,0)
	lst =append(lst, user1.Email, user3.Email)
	lst3 := make([]string,0)
	lst3 =append(lst3, user1.Email, "ngv@gmail.com")
	testCases :=[]struct{
		scenario string
		mockInput r_Request.RequestFriendLists
		Success bool
		Count int
		expectedBody string
		expectedError error
	}{
		{
			scenario: "Success",
			expectedError: nil,
			expectedBody: `pvq@gmail.com`,
			mockInput: r_Request.RequestFriendLists{RequestFriendLists: lst},
			Count: 1,
			Success: true,
		},
		{
			scenario: "Fail by not exists",
			mockInput: r_Request.RequestFriendLists{RequestFriendLists: lst3},
			expectedError: errors.New("no users in table"),
		},
	}
	for _, tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			result, err := repository.NewRepoRelationship().FindCommonListFriend(util.DBInstance,tc.mockInput.RequestFriendLists)
			if tc.scenario == "Success"{
				assert.True(t, tc.Success)
				assert.Equal(t, tc.expectedError, err)
				assert.Equal(t, 1 ,tc.Count)
				assert.Equal(t, tc.expectedBody, result.Friends[0])
			}else{
				assert.Equal(t, tc.expectedError, err)
			}
		})
	}
}
func TestBesubcriber(t *testing.T){
	CreateConnection()
	user1 := &models.User{Email: "hcl@gmail.com"}
	user2 := &models.User{Email: "pvq@gmail.com"}
	user3 := &models.User{Email: "ntt@gmail.com"}
	repository.NewRepo().AddUser(util.DBInstance, user1)
	repository.NewRepo().AddUser(util.DBInstance, user2)
	repository.NewRepo().AddUser(util.DBInstance, user3)
	repository.NewRepoRelationship().AddRelationship(util.DBInstance, user1.Email, user2.Email)
	testCases := []struct{
		scenario string
		mockInput r_Request.RequestUpdate
		expectedError error
		expectedBody bool
	}{
		{
			//update
			scenario: "Success",
			mockInput: r_Request.RequestUpdate{Requestor: user1.Email, Target: user2.Email},
			expectedError: nil,
			expectedBody: true,
		},
		{
			//insert
			scenario: "Success",
			mockInput: r_Request.RequestUpdate{Requestor: user1.Email, Target: user3.Email},
			expectedError: nil,
			expectedBody: true,
		},
		{
			scenario: "Fail by not exists",
			mockInput: r_Request.RequestUpdate{Requestor: user1.Email, Target: "user@gmail.com"},
			expectedError: errors.New("no users in table"),
		},
	}
	for _ , tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			result, err := repository.NewRepoRelationship().BeSubcribe(util.DBInstance,tc.mockInput.Requestor, tc.mockInput.Target)
			if tc.scenario =="Success"{
				assert.Equal(t, tc.expectedError, err)
				assert.Equal(t, tc.expectedBody, result.Success)
			}else{
				assert.Equal(t, tc.expectedError, err)
			}
		})
	}
}
func TestToBlock(t *testing.T){
	CreateConnection()
	user1 := &models.User{Email: "hcl@gmail.com"}
	user2 := &models.User{Email: "pvq@gmail.com"}
	user3 := &models.User{Email: "ntt@gmail.com"}
	repository.NewRepo().AddUser(util.DBInstance, user1)
	repository.NewRepo().AddUser(util.DBInstance, user2)
	repository.NewRepo().AddUser(util.DBInstance, user3)
	repository.NewRepoRelationship().AddRelationship(util.DBInstance, user1.Email, user2.Email)
	testCases := []struct{
		scenario string
		mockInput r_Request.RequestUpdate
		expectedError error
		expectedBody bool
	}{
		{
			//update are friend-->
			scenario: "Success",
			mockInput: r_Request.RequestUpdate{Requestor: user1.Email, Target: user2.Email},
			expectedError: nil,
			expectedBody: true,
		},
		{
			//insert
			scenario: "Success",
			mockInput: r_Request.RequestUpdate{Requestor: user1.Email, Target: user3.Email},
			expectedError: nil,
			expectedBody: true,
		},
		{
			scenario: "Fail by not exists",
			mockInput: r_Request.RequestUpdate{Requestor: user1.Email, Target: "user@gmail.com"},
			expectedError: errors.New("no users in table"),
		},
	}
	for _ , tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			result, err := repository.NewRepoRelationship().ToBlock(util.DBInstance,tc.mockInput.Requestor, tc.mockInput.Target)
			if tc.scenario =="Success"{
				assert.Equal(t, tc.expectedError, err)
				assert.Equal(t, tc.expectedBody, result.Success)
			}else{
				assert.Equal(t, tc.expectedError, err)
			}
		})
	}
}
func TestRetrieveUpdate(t *testing.T){
	CreateConnection()
	user1 := &models.User{Email: "hcl@gmail.com"}
	user2 := &models.User{Email: "pvq@gmail.com"}
	user3 := &models.User{Email: "ntt@gmail.com"}
	repository.NewRepo().AddUser(util.DBInstance, user1)
	repository.NewRepo().AddUser(util.DBInstance, user2)
	repository.NewRepo().AddUser(util.DBInstance,&models.User{Email: "user@gmail.com"})
	repository.NewRepoRelationship().AddRelationship(util.DBInstance, user1.Email, user2.Email)
	testCases := []struct{
		scenario string
		mockInput r_Request.RetrieveUpdate
		expectedError error
		expectedBody bool
		expectedBodyRecipients []string
	}{
		{
			scenario: "Success",
			mockInput: r_Request.RetrieveUpdate{Sender : user1.Email, Tartget: "sent to abc@gmail.com"},
			expectedError: nil,
			expectedBody: true,
			expectedBodyRecipients: []string{
				"pvq@gmail.com",
				"abc@gmail.com",
			},
		},
		{
			scenario: "Fail by not exists",
			mockInput: r_Request.RetrieveUpdate{Sender: user3.Email, Tartget: user1.Email},
			expectedError: errors.New("no users in table"),
		},
	}
	for _ , tc := range testCases{
		t.Run(tc.scenario, func(t *testing.T) {
			result, err := repository.NewRepoRelationship().RetrieveUpdate(util.DBInstance,tc.mockInput.Sender, tc.mockInput.Tartget)
			if tc.scenario =="Success"{
				assert.Equal(t, tc.expectedError, err)
				assert.Equal(t, tc.expectedBody, result.Success)
				assert.Equal(t, tc.expectedBodyRecipients, result.Recipients)
			}else{
				assert.Equal(t, tc.expectedError, err)
			}
		})
	}
}


