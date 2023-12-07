package routers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

// ==========================
// === Структуры для форм ===
// ==========================
type RegisterForm struct {
	FirstName  string `form:"first_name"`
	LastName   string `form:"last_name"`
	MiddleName string `form:"middle_name"`
	Email      string `form:"email"`
	Password   string `form:"password"`
}

// ==========================
// ======== Методы ==========
// ==========================
func register(ctx *fiber.Ctx, db *sql.DB) error {
	formdata := new(RegisterForm)

	if err := ctx.BodyParser(formdata); err != nil {
		return respError(ctx, err.Error())
	}
	// Проверка на существование пользователя
	rows, err := db.Query("SELECT COUNT(*) AS count FROM `users` WHERE 'email' = $1", formdata.Email)
	if err != nil {
		return respError(ctx, err.Error())
	}
	defer rows.Close()
	count, err := checkCount(rows)
	if err != nil {
		return respError(ctx, err.Error())
	}
	if count > 0 {
		return respError(ctx, "Пользователь с таким Email уже существует")
	}

	// Запись в базу
	_, err = db.Query("INSERT INTO `users` (first_name, last_name, middle_name, email, password) VALUES ($1, $2, $3, $4, $5)",
		formdata.FirstName, formdata.LastName, formdata.MiddleName, formdata.Email, formdata.Password)
	if err != nil {
		return respError(ctx, err.Error())
	}

	// Сразу входим
	return nil
}

func checkCount(rows *sql.Rows) (count int, err error) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

func login(ctx *fiber.Ctx, db *sql.DB) error {
	return nil
}
