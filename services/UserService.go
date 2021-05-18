package services

import (
	"Friend_management/models"
	"Friend_management/repository"
	// "database/sql"

	// "database/sql"
	// "errors"

	// "fmt"

	// "fmt"
	"Friend_management/db"
)
type reposervices struct{}
var (
	repo repository.UserRepoInter
)
type UserService interface{
	FindAllUser(database db.Database)(*models.UserList, error)
	AddUser(database db.Database, user *models.User)error
	GetUserByEmail(database db.Database, email string)(models.User, error)
	DeleteUser(database db.Database, email string) error
}
func NewUserService(repository repository.UserRepoInter)UserService{
	repo = repository
	return &reposervices{}
} 
func (* reposervices)FindAllUser(database db.Database)(*models.UserList, error){
	return repo.GetAllUsers(database)
}
func (* reposervices)AddUser(database db.Database, user *models.User)error{
	return repo.AddUser(database,user)
}

func (* reposervices)GetUserByEmail(database db.Database, email string)(models.User, error){
	return repo.GetUserByEmail(database, email)
}
func (* reposervices)DeleteUser(database db.Database, email string) error{
	return repo.DeleteUser(database, email)
}
