package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/roGal1k/golang-beginner/assets/model"
)

func (a *API) createSectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	projectName, ok := vars["projectname"]
	if !ok {
		http.Error(w, "Project name is required"+string(projectName), http.StatusBadRequest)
		return
	}

	idProject := r.Header.Get("Id")
	if idProject == "" {
		http.Error(w, "id project is not found", http.StatusUnauthorized)
		return
	}

	var request model.Section
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.Title == "" {
		http.Error(w, "Section title is required", http.StatusBadRequest)
		return
	}

	num, err := strconv.ParseUint(idProject, 10, 64)
	if err != nil {
		fmt.Println("Ошибка при конвертации:", err)
		return
	}

	request.ProjectID = uint(num)

	// Сохранение секции в базе данных
	result := a.DB.Create(&request)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Ответ пользователю
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "Section created successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (a *API) getSectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Извлечение значений из URL-параметров
	vars := mux.Vars(r)
	projectName, ok := vars["projectname"]
	if !ok {
		http.Error(w, "Project name is required"+string(projectName), http.StatusBadRequest)
		return
	}
	sectionName, ok := vars["sectionname"]
	if !ok {
		http.Error(w, "Section name is required"+string(sectionName), http.StatusBadRequest)
		return
	}

	idProject := r.Header.Get("Id")
	if idProject == "" {
		http.Error(w, "id project is not found", http.StatusUnauthorized)
		return
	}

	var Sections []model.Section
	err := a.DB.Where("ptoject_id = ?", idProject).Find(&Sections).Error
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Sections)
	w.WriteHeader(http.StatusOK)
}

func (a *API) getSectionsHandler(w http.ResponseWriter, r *http.Request) {

}
