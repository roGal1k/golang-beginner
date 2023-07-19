package main

import (
	"log"

	"github.com/roGal1k/golang-beginner/api"
	db "github.com/roGal1k/golang-beginner/internal/database"
)

func main() {
	// Создание экземпляра базы данных
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Создание экземпляра API с передачей базы данных
	api := &api.API{
		DB: database,
	}

	// Запуск сервера
	api.RunServer()
}
