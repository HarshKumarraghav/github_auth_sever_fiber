package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Firstname string             `json:"first_name"`
	Lastname  string             `json:"last_name"`
	Password  string             `json:"password"`
	Email     string             `json:"email"`
	Profile   string             `json:"profile"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Username  string             `json:"username"`
	UserType  string             `json:"user_type"`
}

type InUser struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Profile   string `json:"profile"`
	UserType  string `json:"user_type"`
	Username  string `json:"username"`
}

type OutUser struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Email     string `json:"email"`
	Profile   string `json:"profile"`
	Username  string `json:"username"`
}

func (in *InUser) ToUser() User {
	return User{
		ID:        primitive.NewObjectID(),
		Firstname: in.Firstname,
		Lastname:  in.Lastname,
		Password:  hashPassword(in.Password),
		Email:     in.Email,
		Profile:   in.Profile,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserType:  in.UserType,
		Username:  in.Username,
	}
}

func (u *User) ToOutUser() OutUser {
	return OutUser{
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Email:     u.Email,
		Profile:   u.Profile,
		Username:  u.Username,
	}
}

// hashPassword Changed the raw password string
// into a hashed one which is saved in the database.
func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes)
}
