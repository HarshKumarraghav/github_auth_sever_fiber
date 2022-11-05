package sheets

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sheet struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `json:"Name"`
	Link     string             `json:"Link"`
	Videourl string             `json:"Videourl"`
	Category string             `json:"Category"`
	Level    string             `json:"Level"`
	Problem  string             `json:"problem"`
}
