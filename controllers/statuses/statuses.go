package statuses

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"course-tracker/controllers/errors"
	utils "course-tracker/controllers/utils"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

// ==========================
// === Структуры для форм ===
// ==========================
type StatusesList struct {
	utils.AccessTokenForm
	utils.ListTags
}
type StatusesChangeForm struct {
	Status string `query:"status" form:"status" json:"status"`
	utils.AccessTokenForm
}

// ==================================
// ===  Вспомогательные структуры ===
// ==================================
type StatusInfo struct {
	ApplicationID int64     `json:"application_id"`
	Changer       string    `json:"changer"`
	ChangeDate    time.Time `json:"change_date"`
	Status        int       `json:"status"`
}

// ==========================
// ======== Методы ==========
// ==========================
// ================ /api/statuses/:id ================
func GetStatuses(ctx *fiber.Ctx, db *sql.DB) error {
	// Получение параметров
	form := new(StatusesList)
	if err := ctx.BodyParser(form); err != nil {
		return errors.RespError(ctx, err.Error())
	}

	appid, err := ctx.ParamsInt("id", 0)
	if err != nil {
		return errors.RespError(ctx, "Неверный параметр \"id\"")
	}

	// Перевод в дефолтные значения
	if form.Limit == 0 {
		form.Limit = 10
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
	var result *sql.Row
	var last_id int64
	if user.Perms == 0 {
		result = db.QueryRow("SELECT id FROM course_applications WHERE id = $1 AND student = $2 ORDER BY id DESC LIMIT 1", appid, user.Email)
	} else {
		result = db.QueryRow("SELECT id FROM course_applications WHERE id = $1 ORDER BY id DESC LIMIT 1", appid)
	}
	err = result.Scan(&last_id)
	if err != nil {
		return errors.RespError(ctx, "Ошибка получения данных из БД: "+err.Error())
	}

	if last_id != int64(appid) {
		return errors.RespError(ctx, "Айдишники не совпадают")
	}

	// TODO: Переделать под кастомную фильровку
	rows, err := db.Query("SELECT application_id, changer, change_date, status FROM statuses WHERE application_id = $1 ORDER BY id ASC LIMIT $2 OFFSET $3", last_id, form.Limit, form.Page*form.Limit)
	if err != nil {
		return errors.RespError(ctx, "Ошибка в запросе к БД: "+err.Error())
	}
	defer rows.Close()

	// Создаем массив из записей
	var status_history []any
	for rows.Next() {
		status := &StatusInfo{}
		err := rows.Scan(&status.ApplicationID, &status.Changer, &status.ChangeDate, &status.Status)
		if err != nil {
			fmt.Println(err)
			return errors.RespError(ctx, "Ошибка в формировании списка")
		}
		status_history = append(status_history, status)
	}

	// Формируем JSON ответ
	r := &utils.ResponseStruct{Success: true, Response: status_history}
	response, err := json.Marshal(r)
	if err != nil {
		return errors.RespError(ctx, "Ошибка формирования JSON ответа")
	}

	return ctx.SendString(string(response))
}

// ================ /api/statuses/edit/:id ================
func ChangeStatuses(ctx *fiber.Ctx, db *sql.DB) error {
	// Получение параметров
	form := new(StatusesChangeForm)
	if err := ctx.BodyParser(form); err != nil {
		return errors.RespError(ctx, err.Error())
	}

	appid, err := ctx.ParamsInt("id", 0)
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

	if user.Perms == 0 {
		return errors.RespError(ctx, "Доступ запрещен")
	}

	// Получаем данные из БД
	var result *sql.Row
	var last_id int64
	if user.Perms == 0 {
		result = db.QueryRow("SELECT id FROM course_applications WHERE id = $1 AND student = $2 ORDER BY id DESC LIMIT 1", appid, user.Email)
	} else {
		result = db.QueryRow("SELECT id FROM course_applications WHERE id = $1 ORDER BY id DESC LIMIT 1", appid)
	}
	err = result.Scan(&last_id)
	if err != nil {
		return errors.RespError(ctx, "Ошибка получения данных из БД: "+err.Error())
	}

	if last_id != int64(appid) {
		return errors.RespError(ctx, "Айдишники не совпадают")
	}

	// Ставим статус заявки на "На согласовании"
	_, err = db.Query("INSERT INTO statuses (application_id, changer, change_date, status) VALUES ($1, $2, $3, $4)", last_id, user.Email, time.Now(), form.Status)
	if err != nil {
		return errors.RespError(ctx, "Ошибка записи статуса в БД: "+err.Error())
	}

	return ctx.JSON(&fiber.Map{
		"success": true,
		"response": &fiber.Map{
			"application_id": last_id,
		},
	})
}
