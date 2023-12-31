package auth

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"

	"database/sql"
	"log"

	"course-tracker/controllers/errors"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

var (
	pepper string = "NONE"
)

type Hash struct{}

// ==========================
// === Структуры для форм ===
// ==========================
type RegisterForm struct {
	FirstName  string `form:"first_name" json:"first_name"`
	LastName   string `form:"second_name" json:"second_name"`
	MiddleName string `form:"middle_name" json:"middle_name"`
	Email      string `form:"email" json:"email"`
	Password   string `form:"password" json:"password"`
}

type LoginForm struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func (r *RegisterForm) CheckIsNotNull() bool {
	return r.Email != "" && r.Password != "" && r.FirstName != "" && r.LastName != "" && r.MiddleName != ""
}

func (l *LoginForm) CheckIsNotNull() bool {
	return l.Email != "" && l.Password != ""
}

// ==========================
// ======== Методы ==========
// ==========================
func Register(ctx *fiber.Ctx, db *sql.DB) error {
	formdata := new(RegisterForm)

	if err := ctx.BodyParser(formdata); err != nil {
		return errors.RespError(ctx, err.Error())
	}

	if !formdata.CheckIsNotNull() {
		return errors.RespError(ctx, "Одно из полей пустое")
	}

	// Проверка на существование пользователя
	rows, err := db.Query("SELECT COUNT(*) AS count FROM users WHERE email = $1", formdata.Email)
	if err != nil {
		return errors.RespError(ctx, err.Error())
	}
	defer rows.Close()
	count, err := checkCount(rows)
	if err != nil {
		return errors.RespError(ctx, err.Error())
	}
	if count > 0 {
		return errors.RespError(ctx, "Пользователь с таким Email уже существует")
	}

	// Хэширование пароля
	saltedBytes := []byte(formdata.Password)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return errors.RespError(ctx, "Ошибка с авторизацией на стороне сервера 553")
	}
	hash := string(hashedBytes[:])

	// Запись в базу
	_, err = db.Query("INSERT INTO users (first_name, second_name, middle_name, email, password, perms) VALUES ($1, $2, $3, $4, $5, 0)",
		formdata.FirstName, formdata.LastName, formdata.MiddleName, formdata.Email, hash)
	if err != nil {
		return errors.RespError(ctx, err.Error())
	}

	// Сразу входим
	logindata := new(LoginForm)
	logindata.Email = formdata.Email
	logindata.Password = formdata.Password
	return login(ctx, db, logindata)
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

func LoginFormController(ctx *fiber.Ctx, db *sql.DB) error {
	formdata := new(LoginForm)

	if err := ctx.BodyParser(formdata); err != nil {
		return errors.RespError(ctx, err.Error())
	}

	if !formdata.CheckIsNotNull() {
		return errors.RespError(ctx, "Одно из полей пустое")
	}

	return login(ctx, db, formdata)
}

func login(ctx *fiber.Ctx, db *sql.DB, form *LoginForm) error {
	// Проверка на существование пользователя
	rows, err := db.Query("SELECT COUNT(*) AS count FROM users WHERE email = $1", form.Email)
	if err != nil {
		return errors.RespError(ctx, err.Error())
	}
	defer rows.Close()
	count, err := checkCount(rows)
	if err != nil {
		return errors.RespError(ctx, err.Error())
	}
	if count == 0 {
		return errors.RespError(ctx, "Пользователя с таким Email не существует")
	}

	// Получаем данные и проверяем пароли
	row := db.QueryRow("SELECT email, password FROM users WHERE email = $1", form.Email)

	var email string
	var hashed_pass string
	err = row.Scan(&email, &hashed_pass)
	if err != nil {
		return errors.RespError(ctx, err.Error())
	}

	incoming := []byte(form.Password)
	existing := []byte(hashed_pass)
	err = bcrypt.CompareHashAndPassword(existing, incoming)
	if err != nil {
		return errors.RespError(ctx, "Пароль неверный")
	}

	// Создаем access_token
	timestamp := time.Now().Unix()

	hasher := md5.New()
	_, err = hasher.Write(existing)
	if err != nil {
		return errors.RespError(ctx, err.Error())
	}

	_, err = hasher.Write([]byte(strconv.Itoa(int(timestamp))))
	if err != nil {
		return errors.RespError(ctx, err.Error())
	}

	token := hex.EncodeToString(hasher.Sum(nil))

	_, err = db.Query("INSERT INTO tokens (users_id, access_token) VALUES ($1, $2)", email, token)
	if err != nil {
		return errors.RespError(ctx, err.Error())
	}

	return ctx.JSON(&fiber.Map{
		"success": true,
		"response": &fiber.Map{
			"email":        email,
			"access_token": token,
		},
	})
}
