package routers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

var apiUrl string = "/api/"

func Route(app *fiber.App, db *sql.DB) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("")
	})

	// ===========================
	// =====  Документация =======
	// ===========================
	DocsRoute(app)

	// ===========================
	// ======= Авторизация =======
	// ===========================
	app.Post(apiUrl+"login", func(c *fiber.Ctx) error {
		return register(c, db)
	})
	app.Post(apiUrl+"register", func(c *fiber.Ctx) error {
		return login_form(c, db)
	})

	app.Get(apiUrl+"test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name": "Grame",
			"age": fiber.Map{
				"name": "Grame",
				"age":  20,
			},
		})
	})

	// ===========================
	// ==== Работа с заявками ====
	// ===========================
}

func DocsRoute(app *fiber.App) {
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.SendString("Здесь должна быть ваша документация")
	})
}
