package routers

import "github.com/gofiber/fiber/v2"

var apiUrl string = "/api/"

func Route(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("<h1>Бэкенд сервиса отслеживания обучения</h1><a href=\"/docs\">Документация</a>")
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
