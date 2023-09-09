package model

import "gorm.io/gorm"

// Модель пользователя
type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
	Token    string
	Projects []Project // Связь с проектами пользователя
}

// Модель проекта
type Project struct {
	gorm.Model
	UserID   uint
	Name     string
	Sections []*Section `json:"sections"`
}

// Модель раздела проекта
type Section struct {
	gorm.Model
	ProjectID uint
	Title     string
	//	Image     string
	Contents []Content // Связь с содержимым раздела
}

// Модель содержимого раздела
type Content struct {
	gorm.Model
	SectionID uint
	Type      string
	Data      string
}
