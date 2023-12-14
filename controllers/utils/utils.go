package applications

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

// ==========================
// === Структуры для форм ===
// ==========================
type AccessTokenForm struct {
	AccessToken string `query:"access_token" form:"access_token"`
}

type ListTags struct {
	Limit int `query:"limit" form:"limit"`
	Page  int `query:"page" form:"page"`
}

type Users struct {
	Email      string `form:"email"`
	FirstName  string `form:"first_name"`
	SecondName string `form:"second_name"`
	MiddleName string `form:"middle_name"`
	Perms      int    `form:"perms"`
}

// ==========================
// ======== Методы ==========
// ==========================
func GetUser(token string, db *sql.DB) (user *Users, errorText string) {
	row := db.QueryRow("SELECT email, first_name, second_name, middle_name, perms FROM users JOIN tokens ON tokens.users_id = users.email WHERE access_token = $1", token)

	userdata := new(Users)
	err := row.Scan(&userdata.Email, &userdata.FirstName, &userdata.SecondName, &userdata.MiddleName, &userdata.Perms)
	if err != nil {
		return nil, err.Error()
	}
	if userdata.Email == "" {
		return nil, ""
	}
	return userdata, ""
}

func templateMethod(ctx *fiber.Ctx, db *sql.DB) error {
	return ctx.SendString("Template")
}
