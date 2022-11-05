package routes

import (
	"alphacoder/pkg/user"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/shareed2k/goth_fiber"
)

func githubAuthHandler(repo *user.Repo) fiber.Handler {
	return goth_fiber.BeginAuthHandler
}

func githubAuthCallbackHandler(repo *user.Repo) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userData, err := goth_fiber.CompleteUserAuth(ctx)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error(), "status": "failed"})
		}
		split := strings.Split(userData.Name, " ")
		firstName, lastName := split[0], split[1]

		us, err := repo.ReadByEmail(userData.Email)
		code := ctx.Query("code")
		state := ctx.Query("state")
		claims := jwt.MapClaims{
			"email": userData.Email,
			"code":  code,
			"state": state,
			"admin": true,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		}

		if err != nil || us.Email == "" {

			_, err := repo.Create(user.InUser{
				Firstname: firstName,
				Lastname:  lastName,
				Password:  "rand.Int()",
				Email:     userData.Email,
				Profile:   userData.AvatarURL,
				UserType:  "user",
				Username:  userData.NickName,
			})
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{"error": err.Error(), "status": "failed"})
			}
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error(), "status": "failed"})
		}

		return ctx.Status(200).JSON(fiber.Map{"token": t, "status": "success"})

	}
}

func jwtLoginHandler(repo *user.Repo, svc user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var creds user.AuthBody
		if err := c.BodyParser(&creds); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "status": "sucfailedcess"})
		}

		token, err := svc.Login(creds.Email, creds.Password)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "status": "failed"})
		}

		return c.Status(200).JSON(fiber.Map{"token": token, "status": "success"})
	}
}

func jwtSignupHandler(repo *user.Repo, svc user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in user.InUser
		if err := c.BodyParser(&in); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "status": "sucfailedcess"})
		}

		token, err := svc.SignUp(in)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "status": "failed"})
		}

		return c.Status(200).JSON(fiber.Map{"token": token, "status": "success"})
	}
}

func CreateAuthRoutes(app *fiber.App, userRepo *user.Repo, svc user.Service) {
	app.Post("/login/jwt", jwtLoginHandler(userRepo, svc))
	app.Post("/signup/jwt", jwtSignupHandler(userRepo, svc))
	app.Get("/login/:provider", githubAuthHandler(userRepo))
	app.Get("/auth/callback/github", githubAuthCallbackHandler(userRepo))
}
