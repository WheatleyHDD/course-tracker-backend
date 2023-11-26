package main

import (
	"github.com/gofiber/fiber/v2"

	"course-tracker/routers"
)

func main() {
	app := fiber.New()

	connectDB()
	routers.Route(app)

	app.Listen(":3000")
}
