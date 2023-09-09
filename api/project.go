package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/roGal1k/golang-beginner/assets/model"
)

// Get projects list
func (a *API) getProjectsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims, err := a.authorization(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Получение пользователя из базы данных
	var user model.User
	err = a.DB.Where("username = ?", claims.Username).First(&user).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Получение списка проектов пользователя
	var projects []model.Project
	err = a.DB.Where("user_id = ?", user.ID).Preload("Sections").Find(&projects).Error
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}

	// Отправка списка проектов в ответе
	json.NewEncoder(w).Encode(projects)
}

// Create project
func (a *API) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims, err := a.authorization(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var user model.User
	err = a.DB.Where("username = ?", claims.Username).First(&user).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Создание нового проекта
	var request model.Project
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка наличия обязательных полей
	if request.Name == "" {
		http.Error(w, "Project name is required", http.StatusBadRequest)
		return
	}

	request.UserID = uint(user.ID)
	fmt.Printf("UserID: %d, Username: %s\n", user.ID, claims.Username)

	// Сохранение проекта в базе данных
	result := a.DB.Create(&request)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Ответ пользователю
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "Project created successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// Geted Project
func (a *API) getProjectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId := r.Header.Get("Id")
	if userId == "" {
		http.Error(w, "Authorization token required", http.StatusUnauthorized)
		return
	}
	// Извлечение имени проекта из URL
	projectName := mux.Vars(r)["projectname"]

	// Декодирование имени проекта из URL
	decodedProjectName, err := url.QueryUnescape(projectName)
	if err != nil {
		http.Error(w, "Invalid project name", http.StatusBadRequest)
		return
	}

	// Запрос к базе данных для получения проекта по его имени (или другому идентификатору)
	var project model.Project
	result := a.DB.Where("name = ? AND user_id = ?", decodedProjectName, userId).Find(&project)
	if result.Error != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Ответ пользователю с информацией о проекте
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

func (a *API) updateProjectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Извлечение ID проекта из заголовка
	projectIDStr := r.Header.Get("Id")
	if projectIDStr == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	// Преобразование ID проекта в uint
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// Запрос к базе данных для получения проекта по его ID
	var existingProject model.Project
	result := a.DB.First(&existingProject, uint(projectID))
	if result.Error != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Декодирование JSON-данных с обновленной информацией о проекте
	var updatedProject model.Project
	err = json.NewDecoder(r.Body).Decode(&updatedProject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Выполнение обновления проекта в базе данных
	result = a.DB.Model(&existingProject).Updates(&updatedProject)
	if result.Error != nil {
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	// Ответ пользователю
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Project updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}
