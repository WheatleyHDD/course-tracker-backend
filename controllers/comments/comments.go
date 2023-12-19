package comments

import (
	"course-tracker/controllers/errors"
	utils "course-tracker/controllers/utils"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

// ==========================
// === Структуры для форм ===
// ==========================
type AddCommentForm struct {
	Message string `query:"message" form:"message"`
	utils.AccessTokenForm
}

// ==========================
// ======== Методы ==========
// ==========================

// ================ /api/comments/:application_id ================
func GetComments(ctx *fiber.Ctx, db *sql.DB) error {
	// Получение параметров
	form := new(utils.AccessTokenForm)
	if err := ctx.BodyParser(form); err != nil {
		return errors.RespError(ctx, err.Error())
	}

	appid, err := ctx.ParamsInt("application_id", 0)
	if err != nil {
		return errors.RespError(ctx, "Неверный параметр \"application_id\"")
	}

	// Получение данных пользователя
	user, errs := utils.GetUser(form.AccessToken, db)
	if user == nil {
		if errs == "" {
			return errors.RespError(ctx, "Недействительный access_token")
		}
		return errors.RespError(ctx, errs)
	}

	application, err := utils.GetApplication(appid, user, db)
	if err != nil {
		return errors.RespError(ctx, "Ошибка получения заявки: "+err.Error())
	}

	// Получаем данные из БД
	rows, err := db.Query("SELECT * FROM comms WHERE application_id = $1 ORDER BY id DESC", application.ID)
	if err != nil {
		return errors.RespError(ctx, "Ошибка в запросе к БД: "+err.Error())
	}
	defer rows.Close()

	// Создаем массив из записей
	var comms []*utils.Comment
	for rows.Next() {
		comm := &utils.Comment{}
		err := rows.Scan(&comm.ID, &comm.ApplicationID, &comm.Sender, &comm.Text, &comm.CommsTime)
		if err != nil {
			fmt.Println(err)
			return errors.RespError(ctx, "Ошибка в формировании списка")
		}
		comms = append(comms, comm)
	}

	return ctx.JSON(&fiber.Map{
		"success": true,
		"response": &fiber.Map{
			"application_id": application.ID,
			"comms":          comms,
		},
	})
}

// ================ /api/comments/:application_id/add ================
func AddComment(ctx *fiber.Ctx, db *sql.DB) error {
	// Получение параметров
	form := new(AddCommentForm)
	if err := ctx.BodyParser(form); err != nil {
		return errors.RespError(ctx, err.Error())
	}

	appid, err := ctx.ParamsInt("application_id", 0)
	if err != nil {
		return errors.RespError(ctx, "Неверный параметр \"application_id\"")
	}

	// Получение данных пользователя
	user, errs := utils.GetUser(form.AccessToken, db)
	if user == nil {
		if errs == "" {
			return errors.RespError(ctx, "Недействительный access_token")
		}
		return errors.RespError(ctx, errs)
	}

	application, err := utils.GetApplication(appid, user, db)
	if err != nil {
		return errors.RespError(ctx, "Ошибка получения заявки: "+err.Error())
	}

	// Добавляем комментарий
	_, err = db.Query("INSERT INTO comms (application_id, sender, comm_timestamp, message_text) VALUES ($1, $2, $3, $4)", application.ID, user.Email, time.Now(), form.Message)
	if err != nil {
		return errors.RespError(ctx, "Ошибка записи комментария в БД: "+err.Error())
	}

	// Получаем айди последнего комментария
	var comm_last_id int64
	row := db.QueryRow("SELECT id FROM comms WHERE application_id = $1 ORDER BY id DESC LIMIT 1", appid)
	err = row.Scan(&comm_last_id)
	if err != nil {
		return errors.RespError(ctx, "Ошибка получения комментария из БД: "+err.Error())
	}

	return ctx.JSON(&fiber.Map{
		"success": true,
		"response": &fiber.Map{
			"application_id": application.ID,
			"comm_id":        comm_last_id,
		},
	})
}
