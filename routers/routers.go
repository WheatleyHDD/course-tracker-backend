package routers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"course-tracker/controllers/applications"
	"course-tracker/controllers/auth"
	"course-tracker/controllers/comments"
	"course-tracker/controllers/courses"
	"course-tracker/controllers/statuses"
	"course-tracker/controllers/users"

	_ "github.com/lib/pq"
)

var apiUrl string = "/api"

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
	app.Post(apiUrl+"/login", func(c *fiber.Ctx) error {
		return auth.LoginFormController(c, db)
	})
	app.Post(apiUrl+"/register", func(c *fiber.Ctx) error {
		return auth.Register(c, db)
	})

	// ===========================
	// ==== Работа с заявками ====
	// ===========================
	app.Post(apiUrl+"/user-applications", func(c *fiber.Ctx) error {
		return applications.UserApplications(c, db)
	})
	app.Post(apiUrl+"/user-applications/add", func(c *fiber.Ctx) error {
		return applications.AddApplication(c, db)
	})

	app.Post(apiUrl+"/applications/:id<int>", func(c *fiber.Ctx) error {
		return applications.GetApplication(c, db)
	})
	app.Post(apiUrl+"/applications/edit/:id<int>", func(c *fiber.Ctx) error {
		return applications.EditApplication(c, db)
	})

	app.Post(apiUrl+"/applications", func(c *fiber.Ctx) error {
		return applications.GetApplications(c, db)
	})

	// =============================
	// ==== Работа со статусами ====
	// =============================
	app.Post(apiUrl+"/statuses/:id<int>", func(c *fiber.Ctx) error {
		return statuses.GetStatuses(c, db)
	})
	app.Post(apiUrl+"/statuses/edit/:id<int>", func(c *fiber.Ctx) error {
		return statuses.GetStatuses(c, db)
	})

	// =============================
	// ======== Комментарии ========
	// =============================
	app.Post(apiUrl+"/comments/:application_id<int>", func(c *fiber.Ctx) error {
		return comments.GetComments(c, db)
	})
	app.Post(apiUrl+"/comments/:application_id<int>/add", func(c *fiber.Ctx) error {
		return comments.AddComment(c, db)
	})

	// =============================
	// == Получение пользователя ===
	// =============================
	app.Post(apiUrl+"/profile/:email", func(c *fiber.Ctx) error {
		return users.GetProfile(c, db)
	})

	app.Post(apiUrl+"/profile", func(c *fiber.Ctx) error {
		return users.GetProfile(c, db)
	})

	// =============================
	// ========== Курсы ============
	// =============================
	app.Post(apiUrl+"/courses/names", func(c *fiber.Ctx) error {
		return courses.GetAllCourseNames(c, db)
	})
}

func DocsRoute(app *fiber.App) {
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.SendString("Здесь должна быть ваша документация")
	})
}
