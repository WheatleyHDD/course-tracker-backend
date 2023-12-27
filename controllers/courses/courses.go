package courses

import (
	"database/sql"
	"fmt"

	"course-tracker/controllers/errors"
	utils "course-tracker/controllers/utils"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

// ==========================
// === Структуры для форм ===
// ==========================
type GetAllNamesForm struct {
	Like string `query:"like" form:"like" json:"like"`
	utils.AccessTokenForm
}

type TemplateForm struct {
}

// ==========================
// ======== Методы ==========
// ==========================

// ================ /api/courses/names/:like ================
func GetAllCourseNames(ctx *fiber.Ctx, db *sql.DB) error {
	// Получение параметров
	form := new(GetAllNamesForm)
	if err := ctx.BodyParser(form); err != nil {
		return errors.RespError(ctx, err.Error())
	}

	// Получение данных пользователя
	user, errs := utils.GetUser(form.AccessToken, db)
	if user == nil {
		if errs == "" {
			return errors.RespError(ctx, "Недействительный access_token")
		}
		return errors.RespError(ctx, errs)
	}

	// Получаем данные из БД
	rows, err := db.Query("SELECT course_name FROM course_applications WHERE LOWER(course_name) LIKE $1 GROUP BY course_name ORDER BY course_name ASC LIMIT 10", form.Like+"%")
	if err != nil {
		return errors.RespError(ctx, "Ошибка в запросе к БД: "+err.Error())
	}
	defer rows.Close()

	var courses []string
	for rows.Next() {
		var courseName string
		err := rows.Scan(&courseName)
		if err != nil {
			fmt.Println(err)
			continue
			// return errors.RespError(ctx, "Ошибка в формировании списка")
		}

		courses = append(courses, courseName)
	}

	return ctx.JSON(&fiber.Map{
		"success":  true,
		"response": courses,
	})
}
