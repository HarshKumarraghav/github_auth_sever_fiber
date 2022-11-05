package configuration

import (
	"os"
)

type Config struct {
	MongoURI           string
	GithubClientKey    string
	GithubClientSecret string
	CallbackUrl        string
	Port               string
	JwtSecret          string
}

func FromEnv() Config {
	return Config{
		MongoURI:           os.Getenv("MONGODB_URL"),
		GithubClientKey:    os.Getenv("GITHUB_CLIENT_KEY"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		CallbackUrl:        os.Getenv("CALLBACK_URL"),
		Port:               os.Getenv("PORT"),
		JwtSecret:          os.Getenv("JWT_SECRET"),
	}
}
