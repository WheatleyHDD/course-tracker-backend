package comments

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

// ==========================
// === Структуры для форм ===
// ==========================
type AccessTokenForm struct {
	AccessToken string `form:"access_token"`
}

type TemplateForm struct {
}

// ==========================
// ======== Методы ==========
// ==========================
func TemplateMethod(ctx *fiber.Ctx, db *sql.DB) error {
	return ctx.SendString("Template")
}
