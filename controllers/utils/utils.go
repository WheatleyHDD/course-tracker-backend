package applications

import (
	"database/sql"
	"time"

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

type ResponseStruct struct {
	Success  bool  `json:"success"`
	Response []any `json:"response"`
}

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

type Comment struct {
	ID            int64     `json:"id"`
	ApplicationID int64     `json:"application_id"`
	Sender        string    `json:"sender"`
	CommsTime     time.Time `json:"comms_time"`
	Text          string    `json:"text"`
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

func GetApplication(appid int, user *Users, db *sql.DB) (*Application, error) {
	// Получаем данные из БД
	var result *sql.Row
	application := new(Application)
	if user.Perms == 0 {
		result = db.QueryRow("SELECT * FROM cources_and_statuses WHERE id = $1 AND student = $2 LIMIT 1", appid, user.Email)
	} else {
		result = db.QueryRow("SELECT * FROM cources_and_statuses WHERE id = $1 LIMIT 1", appid)
	}
	err := result.Scan(&application.ID, &application.CourseName, &application.Student, &application.Cost, &application.StartDate, &application.EndDate, &application.Point, &application.Status, &application.Changer, &application.ChangeDate)
	if err != nil {
		return nil, err
	}
	return application, nil
}

func templateMethod(ctx *fiber.Ctx, db *sql.DB) error {
	return ctx.SendString("Template")
}
