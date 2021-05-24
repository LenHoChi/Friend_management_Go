package repository

import (
	"Friend_management/db"
	"Friend_management/models"
	mol "Friend_management/mymodels"
	"context"
	"log"

	// "database/sql"
	// "errors"
	// "os/user"

	"github.com/volatiletech/sqlboiler/v4/boil"
)
type repo struct{
}
func NewRepo() UserRepoInter{
	return &repo{}
}
type UserRepoInter interface{
	GetAllUsers(database db.Database)(*models.UserList, error)
	AddUser(database db.Database, user *models.User) error
	GetUserByEmail(database db.Database, email string) (models.User, error)
	DeleteUser(database db.Database, email string) error
}
func (r *repo)GetAllUsers(database db.Database) (*models.UserList, error) {
	list := &models.UserList{}
	rows, err := mol.Users().All(context.Background(),database.Conn)
	if err != nil {
		return list, err
	}
	for _, v := range rows {
		list.Users = append(list.Users, models.User{Email: v.Email})
	}
	return list, nil
}
func ClearTable (database db.Database){
	database.Conn.Query("delete from users")
}
func (r *repo)AddUser(database db.Database, user *models.User) error {
	p := &mol.User{
		Email: user.Email,
	}
	if err := p.Insert(context.Background(), database.Conn, boil.Infer()); err != nil {
		return err
	}
	return nil
}
func (r *repo)GetUserByEmail(database db.Database, email string) (models.User, error) {
	found, err := mol.FindUser(context.Background(), database.Conn, email)
	if err != nil{
		return models.User{}, err
	}
	return models.User{Email: found.Email}, nil
}

func (r *repo)DeleteUser(database db.Database, email string) error {
	us, _ := mol.FindUser(context.Background(), database.Conn, email)
	log.Println(us.Email)
	_, err := us.Delete(context.Background(), database.Conn)
	if err != nil{
		return err
	}
	return nil
}
