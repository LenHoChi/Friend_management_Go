package services

import(
	"Friend_management/models"
	"Friend_management/repository"
	"Friend_management/db"
	r_Response "Friend_management/models/response"
)
type relationshipServices struct{}
var(
	repoS repository.RelationshipInter
)
type RepositoryService interface{
	GetAllRelationship(database db.Database) (*models.RelationshipList, error)
	FindRelationshipByKey(database db.Database, userEmail string, friendEmail string) (models.Relationship, error)
	AddRelationship(database db.Database, userEmail string, friendEmail string) (*r_Response.ResponseSuccess, error)
	FindListFriend(database db.Database, email string) (*r_Response.ResponseListFriend, error)
	FindCommonListFriend(database db.Database, lstEmail []string) (*r_Response.ResponseListFriend, error)
	BeSubcribe(database db.Database, requestor string, target string) (*r_Response.ResponseSuccess, error)
	ToBlock(database db.Database, requestor string, target string) (*r_Response.ResponseSuccess, error)
	RetrieveUpdate(database db.Database, sender string, target string) (*r_Response.ResponseRetrieve, error)
}
func NewRelationshipService (repo repository.RelationshipInter) RepositoryService{
	repoS = repo
	return &relationshipServices{}
}
func (r *relationshipServices)GetAllRelationship(database db.Database) (*models.RelationshipList, error) {
	return repoS.GetAllRelationship(database)
}

func (r *relationshipServices)FindRelationshipByKey(database db.Database, userEmail string, friendEmail string) (models.Relationship, error) {
	return repoS.FindRelationshipByKey(database, userEmail, friendEmail)
}
func (r *relationshipServices)AddRelationship(database db.Database, userEmail string, friendEmail string) (*r_Response.ResponseSuccess, error) {
	return repoS.AddRelationship(database, userEmail, friendEmail)
}

func (r *relationshipServices)FindListFriend(database db.Database, email string) (*r_Response.ResponseListFriend, error) {
	return repoS.FindListFriend(database, email)
}

func (r *relationshipServices)FindCommonListFriend(database db.Database, lstEmail []string) (*r_Response.ResponseListFriend, error) {
	return repoS.FindCommonListFriend(database, lstEmail)
}

func (r *relationshipServices)BeSubcribe(database db.Database, requestor string, target string) (*r_Response.ResponseSuccess, error) {
	return repoS.BeSubcribe(database,requestor,target)
}

func (r *relationshipServices)ToBlock(database db.Database, requestor string, target string) (*r_Response.ResponseSuccess, error) {
	return repoS.ToBlock(database,requestor,target)
}

func (r *relationshipServices)RetrieveUpdate(database db.Database, sender string, target string) (*r_Response.ResponseRetrieve, error) {
	return repoS.RetrieveUpdate(database,sender,target)
}
