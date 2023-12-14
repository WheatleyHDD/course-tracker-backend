package applications

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"course-tracker/controllers/errors"
	utils "course-tracker/controllers/utils"

	"time"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

// ==========================
// === Структуры для форм ===
// ==========================
type UserApplicationsForm struct {
	StudentEmail string `query:"student_email" form:"student_email"`
	utils.AccessTokenForm
	utils.ListTags
}

// ==================================
// ===  Вспомогательные структуры ===
// ==================================
type TemplateResponse struct{}
type Application struct {
	ID         int64     `json:"id"`
	CourseName string    `json:"course_name"`
	Student    string    `json:"student"`
	Cost       int       `json:"cost"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Point      string    `json:"point"`
	Status     int       `json:"status"`
	Changer    string    `json:"changer"`
	ChangeDate time.Time `json:"change_date"`
}

type ResponseStruct struct {
	Success  bool  `json:"success"`
	Response []any `json:"response"`
}

// ==========================
// ======== Методы ==========
// ==========================

// ================ /api/user-applications ================
func UserApplications(ctx *fiber.Ctx, db *sql.DB) error {
	// Получение параметров
	formdata := new(UserApplicationsForm)

	if err := ctx.BodyParser(formdata); err != nil {
		return errors.RespError(ctx, err.Error())
	}

	// Получение данных пользователя
	user, errs := utils.GetUser(formdata.AccessToken, db)
	if user == nil {
		if errs == "" {
			return errors.RespError(ctx, "Недействительный access_token")
		}
		return errors.RespError(ctx, errs)
	}

	// Перевод в дефолтные значения
	if formdata.Limit == 0 {
		formdata.Limit = 10
	}
	if user.Perms == 0 {
		formdata.StudentEmail = user.Email
	}

	// Защита
	if formdata.StudentEmail == "" {
		return errors.RespError(ctx, "Не указан email пользователя")
	}

	// TODO: Переделать под кастомную фильтрацию
	rows, err := db.Query("SELECT * FROM cources_and_statuses WHERE student = $1 ORDER BY id ASC LIMIT $2 OFFSET $3", formdata.StudentEmail, formdata.Limit, formdata.Page*formdata.Limit)
	if err != nil {
		return errors.RespError(ctx, "Ошибка в запросе к БД")
	}
	defer rows.Close()

	// Создаем массив из записей
	var courses []any
	for rows.Next() {
		course := &Application{}
		err := rows.Scan(&course.ID, &course.CourseName, &course.Student, &course.Cost, &course.StartDate, &course.EndDate, &course.Point, &course.Status, &course.Changer, &course.ChangeDate)
		if err != nil {
			fmt.Println(err)
			return errors.RespError(ctx, "Ошибка в формировании списка")
		}
		courses = append(courses, course)
	}

	// Формируем JSON ответ
	r := &ResponseStruct{true, courses}
	response, err := json.Marshal(r)
	if err != nil {
		return errors.RespError(ctx, "Ошибка формирования JSON ответа")
	}

	return ctx.SendString(string(response))
}
