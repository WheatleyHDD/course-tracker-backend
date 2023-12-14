package main

import (
	"github.com/gofiber/fiber/v2"

	"course-tracker/routers"
)

func main() {
	app := fiber.New()

	ConnectDB()
	routers.Route(app, Db)
	defer Db.Close()

	app.Listen(":3000")
}
