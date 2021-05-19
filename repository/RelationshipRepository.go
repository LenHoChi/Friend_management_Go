package repository

import (
	"Friend_management/db"
	"Friend_management/models"
	r_Response "Friend_management/models/response"
	"database/sql"
	"errors"
	"regexp"
	"strings"
)
type repoRelationship struct{

}
func NewRepoRelationship() RelationshipInter{
	return &repoRelationship{}
}
type RelationshipInter interface{
	GetAllRelationship(database db.Database) (*models.RelationshipList, error)
	FindRelationshipByKey(database db.Database, userEmail string, friendEmail string) (models.Relationship, error)
	AddRelationship(database db.Database, userEmail string, friendEmail string) (*r_Response.ResponseSuccess, error)
	FindListFriend(database db.Database, email string) (*r_Response.ResponseListFriend, error)
	FindCommonListFriend(database db.Database, lstEmail []string) (*r_Response.ResponseListFriend, error)
	BeSubcribe(database db.Database, requestor string, target string) (*r_Response.ResponseSuccess, error)
	ToBlock(database db.Database, requestor string, target string) (*r_Response.ResponseSuccess, error)
	RetrieveUpdate(database db.Database, sender string, target string) (*r_Response.ResponseRetrieve, error)
}
func (r *repoRelationship)GetAllRelationship(database db.Database) (*models.RelationshipList, error) {
	list := &models.RelationshipList{}

	rows, errFind := database.Conn.Query("SELECT * FROM relationship")
	if errFind != nil {
		return list, errFind
	}
	for rows.Next() {
		var relationship models.Relationship
		errScan := rows.Scan(&relationship.UserEmail, &relationship.FriendEmail, &relationship.AreFriend, &relationship.IsSubcriber, &relationship.IsBlock)
		if errScan != nil {
			return list, errScan
		}
		list.Relationships = append(list.Relationships, relationship)
	}
	return list, nil
}

func (r *repoRelationship)FindRelationshipByKey(database db.Database, userEmail string, friendEmail string) (models.Relationship, error) {
	relationship := models.Relationship{}
	query := `select * from relationship where user_email=$1 and friend_email=$2`
	errFind := database.Conn.QueryRow(query, userEmail, friendEmail).Scan(&relationship.UserEmail, &relationship.FriendEmail, &relationship.AreFriend, &relationship.IsSubcriber, &relationship.IsBlock)
	if errFind != nil {
		if errFind == sql.ErrNoRows {
			return relationship, errFind
		}
		return relationship, errFind
	}
	return relationship, nil
}
func (r *repoRelationship)AddRelationship(database db.Database, userEmail string, friendEmail string) (*r_Response.ResponseSuccess, error) {
	//check email similar
	_, errFindUser1 := NewRepo().GetUserByEmail(database, userEmail)
	_, errFindUser2 := NewRepo().GetUserByEmail(database, friendEmail)
	if errFindUser1 != nil || errFindUser2 != nil {
		return nil, errors.New("user not exists")
	}
	//check relationship similar
	//check case have already this relationship but friend is not -->transfer--> true
	_, errFind := r.FindRelationshipByKey(database, userEmail, friendEmail)
	if errFind == nil {
		return nil, errors.New("this relationship exists already")
	}
	//create new relationship
	query := `INSERT INTO relationship values ($1, $2, $3, $4, $5)`
	// database.Conn.QueryRow(query, userEmail, friendEmail, true, false, false)
	_, errInsert := database.Conn.Exec(query, userEmail, friendEmail, true, false, false)
	if errInsert != nil {
		return nil, errors.New("Error: " + errInsert.Error())
	}
	//---
	// database.Conn.QueryRow(query, userEmail, friendEmail, true, false, false)
	_, errInsert2 := database.Conn.Exec(query, friendEmail, userEmail, true, false, false)
	if errInsert2 != nil {
		return nil, errors.New("Error: " + errInsert.Error())
	}
	
	return &r_Response.ResponseSuccess{Success: true}, nil
}

func (r *repoRelationship)FindListFriend(database db.Database, email string) (*r_Response.ResponseListFriend, error) {
	//check emai exists
	_, errFindUser := NewRepo().GetUserByEmail(database, email)
	if errFindUser != nil {
		return nil, errors.New("no users in table")
	}
	list := &r_Response.ResponseListFriend{}
	query := `select friend_email from relationship where user_email = $1 and arefriends = true
	 union
	 select user_email from relationship where friend_email = $1 and arefriends = true`

	rows, errFindFriend := database.Conn.Query(query, email)

	if errFindFriend != nil {
		return list, errFindFriend
	}
	for rows.Next() {
		var email string
		errScan := rows.Scan(&email)
		if errScan != nil {
			return nil, errScan
		}
		list.Friends = append(list.Friends, email)
	}
	list.Success = true
	list.Count = len(list.Friends)
	return list, nil
}

func (r *repoRelationship)FindCommonListFriend(database db.Database, lstEmail []string) (*r_Response.ResponseListFriend, error) {
	list := &r_Response.ResponseListFriend{}
	//check same email
	//check exists email
	_, errFindUser1 := NewRepo().GetUserByEmail(database, lstEmail[0])
	_, errFindUser2 := NewRepo().GetUserByEmail(database, lstEmail[1])
	if errFindUser1 != nil || errFindUser2 != nil {
		return nil, errors.New("no users in table")
	}
	query := `select r.friend_email from relationship r
	where r.user_email in ($1,$2) and r.arefriends =true 
	group by r.friend_email 
	having count(r.friend_email)>1`
	rows, errFindFriend := database.Conn.Query(query, lstEmail[0], lstEmail[1])

	if errFindFriend != nil {
		return list, errFindFriend
	}
	for rows.Next() {
		var email string
		errScan := rows.Scan(&email)
		if errScan != nil {
			return nil, errScan
		}
		list.Friends = append(list.Friends, email)
	}
	list.Success = true
	list.Count = len(list.Friends)
	return list, nil
}

func (r *repoRelationship)BeSubcribe(database db.Database, requestor string, target string) (*r_Response.ResponseSuccess, error) {
	//check case have already this relationship but issbucriber is not -->transfer--> true
	queryUpdate := `update relationship set issubcriber =true where user_email =$1 and friend_email =$2`
	queryInsert := `INSERT INTO relationship values ($1, $2, $3, $4, $5)`
	// database.Conn.QueryRow(query, requestor, target)
	//check exists email
	_, errFindUser1 := NewRepo().GetUserByEmail(database, requestor)
	_, errFindUser2 := NewRepo().GetUserByEmail(database, target)
	if errFindUser1 != nil || errFindUser2 != nil {
		return nil, errors.New("no users in table")
	}
	_, errFindRelationship := r.FindRelationshipByKey(database, requestor, target)
	//not exitst-->create new
	if errFindRelationship != nil {
		_, errInsert := database.Conn.Exec(queryInsert, requestor, target, false, true, false)
		if errInsert != nil {
			return nil, errInsert
		}
	} else {
		_, errUpdate := database.Conn.Exec(queryUpdate, requestor, target)
		if errUpdate != nil {
			return nil, errUpdate
		}
	}
	return &r_Response.ResponseSuccess{Success: true}, nil
}

func (r *repoRelationship)ToBlock(database db.Database, requestor string, target string) (*r_Response.ResponseSuccess, error) {
	queryInsert := `INSERT INTO relationship values ($1, $2, $3, $4, $5)`
	queryUpdate := `update relationship set issubcriber =false where user_email=$1 and friend_email=$2`
	queryUpdateBlock := `update relationship set issubcriber =false , isblock=true where user_email=$1 and friend_email=$2`
	_, errFindUser1 := NewRepo().GetUserByEmail(database, requestor)
	_, errFindUser2 := NewRepo().GetUserByEmail(database, target)
	if errFindUser1 != nil || errFindUser2 != nil {
		return nil, errors.New("no users in table")
	}
	re, errFindRelationship := r.FindRelationshipByKey(database, requestor, target)
	if errFindRelationship != nil {
		_, errInsert := database.Conn.Exec(queryInsert, requestor, target, false, false, true)
		if errInsert != nil {
			return nil, errInsert
		}
	} else {
		if !re.AreFriend {
			_, errUpdateBlock := database.Conn.Exec(queryUpdateBlock, requestor, target)
			if errUpdateBlock != nil {
				return nil, errUpdateBlock
			}
		} else {
			_, errUpdate := database.Conn.Exec(queryUpdate, requestor, target)
			if errUpdate != nil {
				return nil, errUpdate
			}
		}
	}
	return &r_Response.ResponseSuccess{Success: true}, nil
}

func (r *repoRelationship)RetrieveUpdate(database db.Database, sender string, target string) (*r_Response.ResponseRetrieve, error) {
	_, errFindUser := NewRepo().GetUserByEmail(database, sender)
	if errFindUser != nil {
		return nil, errors.New("no users in table")
	}
	list := &r_Response.ResponseRetrieve{}
	query := `select friend_email from relationship 
	where user_email =$1 and (arefriends=true or issubcriber=true)
	and isblock =false`
	rows, errFindFriend := database.Conn.Query(query, sender)
	if errFindFriend != nil {
		return list, errFindFriend
	}
	for rows.Next() {
		var email string
		errScan := rows.Scan(&email)
		if errScan != nil {
			return nil, errScan
		}
		list.Recipients = append(list.Recipients, email)
	}
	lstTemp := CheckString(target)
	for _, i := range lstTemp {
		if IsEmailValid(i) {
			list.Recipients = append(list.Recipients, i)
		}
	}
	list.Success = true
	return list, nil
}
func CheckString(text string) []string {
	split := strings.Split(text, " ")
	lstEmail := make([]string, 0)
	for _, i := range split {
		if CheckContain(i) {
			lstEmail = append(lstEmail, i)
		}
	}
	return lstEmail
}
func CheckContain(str string) bool {
	bool := strings.Contains(str, "@")
	return bool
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !emailRegex.MatchString(e) {
		return false
	}
	return true
}
