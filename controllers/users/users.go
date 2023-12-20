package users

import (
	"database/sql"

	"course-tracker/controllers/errors"
	utils "course-tracker/controllers/utils"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

// ==========================
// === Структуры для форм ===
// ==========================

// ==========================
// ======== Методы ==========
// ==========================
func GetProfile(ctx *fiber.Ctx, db *sql.DB) error {
	// Получение параметров
	form := new(utils.AccessTokenForm)
	if err := ctx.BodyParser(form); err != nil {
		return errors.RespError(ctx, err.Error())
	}

	email := ctx.Params("email", "")

	// Если нет емейла - возвращаем текущего пользователя
	if email == "" {
		// Получение данных пользователя
		user, errs := utils.GetUser(form.AccessToken, db)
		if user == nil {
			if errs == "" {
				return errors.RespError(ctx, "Недействительный access_token")
			}
			// ХУЙня не работает
			// Почини, хуесос
			// Говно уебищное
			return errors.RespError(ctx, errs)
		}

		return ctx.JSON(&fiber.Map{
			"success": true,
			"response": &fiber.Map{
				"email":       user.Email,
				"first_name":  user.FirstName,
				"second_name": user.SecondName,
				"middle_name": user.MiddleName,
				"perms":       user.Perms,
				"token":       form.AccessToken,
			},
		})
	}

	// возвращаем необходимого
	user, err := getProfile(db, email)
	if err != nil {
		return errors.RespError(ctx, "Ошибка получения пользователя: "+err.Error())
	}

	return ctx.JSON(&fiber.Map{
		"success": true,
		"response": &fiber.Map{
			"email":       user.Email,
			"first_name":  user.FirstName,
			"second_name": user.SecondName,
			"middle_name": user.MiddleName,
			"perms":       user.Perms,
		},
	})
}

func getProfile(db *sql.DB, email string) (*utils.Users, error) {
	result := new(utils.Users)

	row := db.QueryRow("SELECT * FROM user_info WHERE email = $1", email)
	err := row.Scan(&result.Email, &result.FirstName, &result.SecondName, &result.MiddleName, &result.Perms)
	if err != nil {
		return nil, err
	}

	return result, nil
}
