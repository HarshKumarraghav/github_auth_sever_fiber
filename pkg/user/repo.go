package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository is an interfaces that defines the schema of
// the CRUD operations that can be performed on the User
// entity. The Implementation might be changed later in
// case we migrate away from gorm.
type Repository interface {
	Create(in InUser) (User, error)
	Read(id string) (User, error)
	Update(id string, upd map[string]interface{}) (User, error)
	Delete(string int) bool
	ReadByEmail(email string) (User, error)
	ReadByID(id string) (User, error)
	ReadByUsernanme(username string) (User, error)
}

// Repo is the struct that Implements the Repository Interface.
// To Create a Repo, Use the NewRepo Function, it takes in a DB of type *gorm.DB
type Repo struct {
	db      *mongo.Collection
	context context.Context
}

// ReadByEmail is used to fetch a user from the database with their email.
func (s *Repo) ReadByEmail(email string) (User, error) {
	var user User
	err := s.db.FindOne(s.context, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return user, errors.New("user not found with this email")
	}
	return user, nil
}

// ReadByEmail is used to fetch a user from the database with their email.
func (s *Repo) ReadByUsernanme(username string) (User, error) {
	var user User
	err := s.db.FindOne(s.context, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return user, errors.New("user not found with this email")
	}
	return user, nil
}

// ReadByUsername is used to fetch a user from the database with their username.
func (s *Repo) ReadByID(id string) (User, error) {
	var user User
	err := s.db.FindOne(s.context, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return user, errors.New("user not found with this username")
	}
	return user, nil
}

// Create adds a User entity into the database.
// It takes in a InUser entity which is later turned
// into a User struct.
func (s *Repo) Create(in InUser) (User, error) {
	user := in.ToUser()
	_, err := s.db.InsertOne(s.context, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// Read fetches a User entity into the database.
// :param: id: int -> User.ID
func (s *Repo) Read(id string) (User, error) {
	var user User
	err := s.db.FindOne(s.context, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return user, errors.New("user not found with this id")
	}
	return user, nil
}

// Update updates a User entity already present in the database.
// :param: id: int -> User.ID
// :param: upd: map[string]interface{} -> Weakly Typed Parameters From InUser struct
func (s *Repo) Update(id string, upd map[string]interface{}) (User, error) {
	var u User
	if err := s.db.FindOneAndUpdate(s.context, bson.M{"_id": id}, upd).Decode(&u); err != nil {
		return u, err
	}
	return u, nil
}

// Delete deletes a user from the database
func (s *Repo) Delete(id int) bool {
	delete, err := s.db.DeleteOne(s.context, bson.M{"_id": id})
	if err != nil {
		return false
	}
	return delete.DeletedCount == 1
}

// NewRepo returns a Repo which can be used for various operations later.
func NewRepo(db *mongo.Database) Repository {
	ctx := context.TODO()
	return &Repo{db: db.Collection("users"), context: ctx}
}
