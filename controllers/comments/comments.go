package comments

import (
	_ "github.com/lib/pq"
)

// ==========================
// === Структуры для форм ===
// ==========================
type TemplateForm struct {
}

// ==========================
// ======== Методы ==========
// ==========================

// ================ /api/comments/:application_id ================
/*
func GetComments(ctx *fiber.Ctx, db *sql.DB) error {
	// Получение параметров
	form := new(utils.AccessTokenForm)
	if err := ctx.BodyParser(form); err != nil {
		return errors.RespError(ctx, err.Error())
	}

	appid, err := ctx.ParamsInt("application_id", 0)
	if err != nil {
		return errors.RespError(ctx, "Неверный параметр \"id\"")
	}

	// Получение данных пользователя
	user, errs := utils.GetUser(form.AccessToken, db)
	if user == nil {
		if errs == "" {
			return errors.RespError(ctx, "Недействительный access_token")
		}
		return errors.RespError(ctx, errs)
	}
}
*/
