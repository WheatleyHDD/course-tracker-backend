package routers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"course-tracker/controllers/applications"
	"course-tracker/controllers/auth"

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
		return auth.LoginFormController(c, db)
	})
	app.Post(apiUrl+"register", func(c *fiber.Ctx) error {
		return auth.Register(c, db)
	})

	// ===========================
	// ==== Работа с заявками ====
	// ===========================
	app.Post(apiUrl+"user-applications", func(c *fiber.Ctx) error {
		return applications.UserApplications(c, db)
	})
	app.Post(apiUrl+"user-applications/add", func(c *fiber.Ctx) error {
		return applications.AddApplication(c, db)
	})

	app.Post(apiUrl+"applications/:id", func(c *fiber.Ctx) error {
		return applications.GetApplication(c, db)
	})
}

func DocsRoute(app *fiber.App) {
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.SendString("Здесь должна быть ваша документация")
	})
}
