package routers

import "github.com/gofiber/fiber/v2"

func respError(ctx *fiber.Ctx, err string) error {
	return ctx.JSON(fiber.Map{
		"success": false,
		"error":   err,
	})
}
