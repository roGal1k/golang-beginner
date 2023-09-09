package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/roGal1k/golang-beginner/assets/model"
)

type API struct {
	DB *gorm.DB
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Функция для верификации токена и получения информации о пользователе
func (a *API) authorization(r *http.Request) (*Claims, error) {
	// Извлекаем токен из запроса
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("authorization token required")
	}

	// Верификация токена
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil // Замените "your-secret-key" на ваш секретный ключ
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	// Проверка валидности токена и извлечение информации о пользователе
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (a *API) registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка наличия обязательных полей
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	expirationTime := time.Now().Add(31 * 24 * time.Hour) // Срок действия токена - 31 день
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	fmt.Print(claims.Id)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("your-secret-key")) // Замените "your-secret-key" на ваш секретный ключ
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Сохранение пользователя в базе данных
	err = a.DB.Create(&model.User{
		Username: user.Username,
		Password: string(hashedPassword),
		Token:    signedToken,
	}).Error
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка наличия обязательных полей
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Получение пользователя из базы данных
	var storedUser model.User
	result := a.DB.Where("username = ?", user.Username).First(&storedUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Токен уже должен быть в storedUser.Token после регистрации
	if storedUser.Token == "" {
		http.Error(w, "User token not found", http.StatusUnauthorized)
		return
	}

	// Отправка токена в ответе
	response := map[string]string{
		"token": storedUser.Token,
	}
	fmt.Print(storedUser.ID)

	json.NewEncoder(w).Encode(response)
}

func (a *API) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Реализуйте здесь логику выхода пользователя из системы,
	// например, аннулирование или удаление токена
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"message": "Logged out",
	}
	json.NewEncoder(w).Encode(response)
}

func (a *API) RunServer() {
	// Создание маршрутизатора mux
	r := mux.NewRouter()

	// Настройка маршрутов с использованием mux
	r.HandleFunc("/register", a.registerHandler).Methods("POST")
	r.HandleFunc("/login", a.loginHandler).Methods("POST")
	r.HandleFunc("/logout", a.logoutHandler).Methods("POST")
	r.HandleFunc("/projects", a.getProjectsHandler).Methods("GET")
	r.HandleFunc("/projects/create", a.createProjectHandler).Methods("POST")
	r.HandleFunc("/project/{projectname}/section/create", a.createSectionHandler).Methods("POST")
	r.HandleFunc("/project/{projectname}", a.getProjectHandler).Methods("GET")
	r.HandleFunc("/project/{projectname}/sections", a.getSectionsHandler).Methods("GET")
	r.HandleFunc("/project/{projectname}/section/{sectionname}", a.getSectionHandler).Methods("GET")
	r.HandleFunc("/project/{projectname}/update", a.updateProjectHandler).Methods("PUT")

	// Используйте mux вместо http.ListenAndServe
	http.Handle("/", r)

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
