package sheets

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository is an interfaces that defines the schema of
// the CRUD operations that can be performed on the User
// entity. The Implementation might be changed later in
// case we migrate away from gorm.
type Repository interface {
	ReadCount(count, page int) ([]Sheet, error)
}

// Repo is the struct that Implements the Repository Interface.
// To Create a Repo, Use the NewRepo Function, it takes in a DB of type *gorm.DB
type Repo struct {
	db      *mongo.Collection
	context context.Context
}

// ReadByEmail is used to fetch a Sheet from the database with their email.
func (s *Repo) ReadCount(count, skip int) ([]Sheet, error) {
	var sheets []Sheet
	cursor, err := s.db.Find(s.context, bson.M{}, options.Find().SetLimit(int64(count)).SetSkip(int64(skip)))
	if err != nil {
		return sheets, err
	}
	var sheet Sheet
	for cursor.Next(s.context) {
		cursor.Decode(&sheet)
		sheets = append(sheets, sheet)
	}

	return sheets, nil
}

// NewRepo returns a Repo which can be used for various operations later.
func NewRepo(db *mongo.Database) Repository {
	ctx := context.TODO()
	return &Repo{db: db.Collection("sheets"), context: ctx}
}
