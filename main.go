package main

import (
	"alphacoder/api/routes"
	"alphacoder/pkg/configuration"
	"alphacoder/pkg/sheets"
	"alphacoder/pkg/user"
	"context"
	"log"
	"os"

	jwtware "github.com/gofiber/jwt/v3"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	app := fiber.New()

	def := cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}

	app.Use(cors.New(def))

	godotenv.Load()
	config := configuration.FromEnv()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("alphacoder")

	userRepo := user.NewRepo(db)
	sheetsRepo := sheets.NewRepo(db)

	userSvc := user.NewAuthService(userRepo.(*user.Repo))

	goth.UseProviders(github.New(config.GithubClientKey, config.GithubClientSecret, config.CallbackUrl))

	routes.CreateAuthRoutes(app, userRepo.(*user.Repo), userSvc)
	routes.CreateSheetRoutes(app, sheetsRepo.(*sheets.Repo))

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(config.JwtSecret),
	}))
	app.Listen(":" + os.Getenv("PORT"))
}
