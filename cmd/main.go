package main

import (
	"log"

	"github.com/roGal1k/golang-beginner/api"
	db "github.com/roGal1k/golang-beginner/internal/database"
	"gorm.io/gorm"
)

type API struct {
	DB *gorm.DB
}

func main() {
	// Создание экземпляра базы данных
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// Выполнение автомиграций
	err = db.AutoMigrate(database)
	if err != nil {
		log.Fatal(err)
	}

	// Создание экземпляра API с передачей базы данных
	apiInstance := &api.API{
		DB: database,
	}

	// Запуск сервера
	apiInstance.RunServer()
}
