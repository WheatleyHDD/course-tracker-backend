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
		return c.SendString("Hello, World!")
	})
	app.Post(apiUrl+"register", func(c *fiber.Ctx) error {
		// данные которые получаем
		//
		return c.SendString("Hello, World!")
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
