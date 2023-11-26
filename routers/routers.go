package routers

import "github.com/gofiber/fiber/v2"

var apiUrl string = "/api/"

func Route(app *fiber.App) {
	// ===========================
	// ======= Авторизация =======
	// ===========================
	app.Post(apiUrl+"login", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Post(apiUrl+"register", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// ===========================
	// ==== Работа с заявками ====
	// ===========================
}
