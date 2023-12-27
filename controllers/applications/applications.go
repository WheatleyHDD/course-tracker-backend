package applications

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

	"course-tracker/controllers/errors"
	utils "course-tracker/controllers/utils"

	"time"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

// ==========================
// === Структуры для форм ===
// ==========================
type GetAllApplicationsForm struct {
	utils.AccessTokenForm
	utils.ListTags
}

type UserApplicationsForm struct {
	StudentEmail string `query:"student_email" form:"student_email" json:"student_email"`
	utils.AccessTokenForm
	utils.ListTags
}

type AddApplicationForm struct {
	StudentEmail string `query:"student_email" form:"student_email" json:"student_email"`
	CourseName   string `query:"course_name" form:"course_name" json:"course_name"`
	Cost         int    `query:"cost" form:"cost" json:"cost"`
	StartDate    string `query:"start_date" form:"start_date" json:"start_date"`
	EndDate      string `query:"end_date" form:"end_date" json:"end_date"`
	Point        string `query:"point" form:"point" json:"point"`
	utils.AccessTokenForm
}

// ==================================
// ===  Вспомогательные структуры ===
// ==================================
type ApplicationData struct {
	course_name string
	cost        int
	start_date  string
	end_date    string
	point       string
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
		course := &utils.Application{}
		err := rows.Scan(&course.ID, &course.CourseName, &course.Student, &course.Cost, &course.StartDate, &course.EndDate, &course.Point, &course.Status, &course.Changer, &course.ChangeDate)
		if err != nil {
			fmt.Println(err)
			continue
			// return errors.RespError(ctx, "Ошибка в формировании списка")
		}
		courses = append(courses, course)
	}

	// Формируем JSON ответ
	r := &utils.ResponseStruct{Success: true, Response: courses}
	response, err := json.Marshal(r)
	if err != nil {
		return errors.RespError(ctx, "Ошибка формирования JSON ответа")
	}

	return ctx.SendString(string(response))
}

// ================ /api/user-applications/add ================
func AddApplication(ctx *fiber.Ctx, db *sql.DB) error {
	// Получение параметров
	form := new(AddApplicationForm)
	if err := ctx.BodyParser(form); err != nil {
		return errors.RespError(ctx, err.Error())
	}

	// Проверка на валидность
	if form.CourseName == "" || form.Cost == 0 || form.StartDate == "" || form.EndDate == "" ||
		form.Point == "" || form.AccessToken == "" {
		return errors.RespError(ctx, "Одно или несколько полей незаполнено")
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
		form.StudentEmail = user.Email
	} else if form.StudentEmail == "" {
		return errors.RespError(ctx, "Нет E-mail студента")
	}

	// Конвертируем даты
	start_date, err := time.Parse("2006-01-02", form.StartDate)
	if err != nil {
		return errors.RespError(ctx, "Ошибка конвертации даты")
	}
	end_date, err := time.Parse("2006-01-02", form.EndDate)
	if err != nil {
		return errors.RespError(ctx, "Ошибка конвертации даты")
	}
	if end_date.Sub(start_date) <= 0 {
		return errors.RespError(ctx, "Неверный временной диапазон")
	}

	// Запись в базу
	_, err = db.Query("INSERT INTO course_applications (course_name, student, cost, start_date, end_date, point) VALUES ($1, $2, $3, $4, $5, $6)",
		form.CourseName, form.StudentEmail, form.Cost, start_date, end_date, form.Point)
	if err != nil {
		return errors.RespError(ctx, "Ошибка записи в БД: "+err.Error())
	}

	// Получаем айдишник последней заявки
	var last_id int64
	row := db.QueryRow("SELECT id FROM course_applications WHERE student = $1 ORDER BY id DESC LIMIT 1", form.StudentEmail)
	err = row.Scan(&last_id)
	if err != nil {
		return errors.RespError(ctx, "Ошибка получения последнего айдишника: "+err.Error())
	}

	// Ставим статус заявки на "На согласовании"
	_, err = db.Query("INSERT INTO statuses (application_id, changer, change_date, status) VALUES ($1, $2, $3, $4)", last_id, user.Email, time.Now(), 0)
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

// ================ /api/applications/:id ================
func GetApplication(ctx *fiber.Ctx, db *sql.DB) error {
	// Получение параметров
	form := new(utils.AccessTokenForm)
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

	application, err := utils.GetApplication(appid, user, db)
	if err != nil {
		return errors.RespError(ctx, "Ошибка получения заявок: "+err.Error())
	}

	return ctx.JSON(&fiber.Map{
		"success":  true,
		"response": application,
	})
}

// ================ /api/applications/edit/:id ================
func EditApplication(ctx *fiber.Ctx, db *sql.DB) error {
	// Получение параметров
	form := new(AddApplicationForm)
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

	appid, err := ctx.ParamsInt("id", 0)
	if err != nil {
		return errors.RespError(ctx, "Неверный параметр \"id\"")
	}

	// Можно ли менять заявку?
	if user.Perms == 0 {
		result := db.QueryRow("SELECT id FROM cources_and_statuses WHERE id = $1 AND student = $2 LIMIT 1", appid, user.Email)
		err = result.Scan(&appid)
		if err != nil {
			return errors.RespError(ctx, "Ошибка получения данных из БД: "+err.Error())
		}
	}

	// Ищем измененные поля
	custom := ApplicationData{
		course_name: form.CourseName,
		cost:        form.Cost,
		start_date:  form.StartDate,
		end_date:    form.EndDate,
		point:       form.Point,
	}

	value := reflect.ValueOf(custom)
	numFields := value.NumField()
	structType := value.Type()

	for i := 0; i < numFields; i++ {
		fieldValue := value.Field(i)

		if fieldValue.IsZero() {
			continue
		}

		sql := fmt.Sprintf("UPDATE course_applications SET %s = $1 WHERE id = $2", structType.Field(i).Name)

		if structType.Field(i).Type.Name() == "int" {
			_, err = db.Query(sql, fieldValue.Int(), appid)
		} else {
			_, err = db.Query(sql, fieldValue.String(), appid)
		}
		if err != nil {
			fmt.Println(err)
			return errors.RespError(ctx, "Проблема с обновлением параметра \""+structType.Field(i).Name+"\"")
		}
	}

	return ctx.JSON(&fiber.Map{
		"success": true,
		"response": &fiber.Map{
			"application_id": appid,
		},
	})
}

// ================ /api/applications ================
func GetApplications(ctx *fiber.Ctx, db *sql.DB) error {
	// Получение параметров
	form := new(GetAllApplicationsForm)
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

	if user.Perms == 0 {
		return errors.RespError(ctx, "Нет доступа")
	}

	// Перевод в дефолтные значения
	if form.Limit == 0 {
		form.Limit = 10
	}

	// Получаем данные из БД
	rows, err := db.Query("SELECT * FROM cources_and_statuses ORDER BY id ASC LIMIT $1 OFFSET $2", form.Limit, form.Page*form.Limit)
	if err != nil {
		return errors.RespError(ctx, "Ошибка в запросе к БД: "+err.Error())
	}
	defer rows.Close()

	// Создаем массив из записей
	var courses []any
	for rows.Next() {
		course := &utils.Application{}
		err := rows.Scan(&course.ID, &course.CourseName, &course.Student, &course.Cost, &course.StartDate, &course.EndDate, &course.Point, &course.Status, &course.Changer, &course.ChangeDate)
		if err != nil {
			fmt.Println(err)
			continue
			// return errors.RespError(ctx, "Ошибка в формировании списка")
		}
		courses = append(courses, course)
	}

	// Формируем JSON ответ
	r := &utils.ResponseStruct{Success: true, Response: courses}
	response, err := json.Marshal(r)
	if err != nil {
		return errors.RespError(ctx, "Ошибка формирования JSON ответа")
	}

	return ctx.SendString(string(response))
}
