package user

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(email, password string) (string, error)
	SignUp(in InUser) (string, error)
}

type Svc struct {
	repo *Repo
}

// Login implements Service
func (s *Svc) Login(email string, password string) (string, error) {

	user, err := s.repo.ReadByEmail(email)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"email": email,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return t, nil

}

// SignUp implements Service
func (s *Svc) SignUp(in InUser) (string, error) {
	user, err := s.repo.ReadByEmail(in.Email)
	if err != nil {
		return "", err
	}

	if user.Email == in.Email {
		return "", errors.New("user with id already exists")
	}

	create, err := s.repo.Create(in)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"email": create.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return t, nil
}

func NewAuthService(repo *Repo) Service {
	return &Svc{
		repo: repo,
	}
}
