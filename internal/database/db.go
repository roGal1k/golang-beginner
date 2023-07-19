// Модуль db
package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/roGal1k/golang-beginner/assets/model"

	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	Conn *sql.DB
}

func NewDatabase() (*Database, error) {
	// Инициализация подключения к базе данных
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=admin dbname=mydatabase sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Проверка соединения с базой данных
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established")

	return &Database{
		Conn: db,
	}, nil
}

func (d *Database) SaveMessage(message *model.Message) error {
	// Подготовка SQL-запроса для вставки данных
	query := "INSERT INTO messages (text) VALUES ($1)"
	_, err := d.Conn.Exec(query, message.Text)
	if err != nil {
		return err
	}

	log.Println("Message saved successfully")

	return nil
}

func (d *Database) SaveUser(user *model.User) error {
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"
	_, err := d.Conn.Exec(query, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetUserByUsername(username string) (*model.User, error) {
	query := "SELECT username, password FROM users WHERE username = $1"
	row := d.Conn.QueryRow(query, username)

	var user model.User
	err := row.Scan(&user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (d *Database) RegisterUser(user *model.User) error {
	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Подготовка SQL-запроса для вставки данных пользователя
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"
	_, err = d.Conn.Exec(query, user.Username, string(hashedPassword))
	if err != nil {
		return err
	}

	log.Println("User registered successfully")

	return nil
}
