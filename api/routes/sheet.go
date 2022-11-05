package routes

import (
	"alphacoder/pkg/sheets"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func allSheetsHandler(repo *sheets.Repo) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, err := strconv.Atoi(c.Query("page", "1"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "please enter a numeric value as page query paramter "})
		}
		var skip int
		if page == 1 {
			skip = 0
		} else {
			skip = page * 100
		}
		sheets, err := repo.ReadCount(100, skip)
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": err.Error()})
		}
		return c.Status(200).JSON(fiber.Map{
			"data": sheets,
			"next": fmt.Sprintf("/sheets/all?page=%d", page+1),
		})
	}
}

func CreateSheetRoutes(app *fiber.App, sheetRepo *sheets.Repo) {

	app.Get("/sheets/all", allSheetsHandler(sheetRepo))
}
