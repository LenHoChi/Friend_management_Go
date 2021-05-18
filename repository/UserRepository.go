package repository

import (
	"Friend_management/models"
	"database/sql"
	"errors"
	"Friend_management/db"
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

	rows, err := database.Conn.Query("SELECT * FROM users")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Email)
		if err != nil {
			return list, err
		}
		list.Users = append(list.Users, user)
	}
	return list, nil
}
func ClearTable (database db.Database){
	database.Conn.Query("delete from users")
}
func (r *repo)AddUser(database db.Database, user *models.User) error {

	query := `INSERT INTO users (email) VALUES ($1)`
	_, errFind := r.GetUserByEmail(database, user.Email)
	if errFind == nil {
		return errors.New("this email exists already")
	}
	_, err := database.Conn.Exec(query, user.Email)
	if err != nil {
		return err
	}
	return nil
}
func (r *repo)GetUserByEmail(database db.Database, email string) (models.User, error) {
	user := models.User{}
	query := `select * from users where email = $1;`

	err := database.Conn.QueryRow(query, email).Scan(&user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
		return user, err
	}
	return user, nil
}

func (r *repo)DeleteUser(database db.Database, email string) error {
	if _,err := r.GetUserByEmail(database, email); err!= nil{
		return errors.New("this user not exists")
	}
	query := `delete from users where email =$1`
	_, err := database.Conn.Exec(query, email)
	switch err {
	case sql.ErrNoRows:
		return db.ErrNoMatch
	default:
		return nil
	}
}
