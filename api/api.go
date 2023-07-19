// Модуль api
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/roGal1k/golang-beginner/assets/model"
	db "github.com/roGal1k/golang-beginner/internal/database"
)

type API struct {
	DB *db.Database
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
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

	// Сохранение пользователя в базе данных
	err = a.DB.SaveUser(&model.User{
		Username: user.Username,
		Password: string(hashedPassword),
	})
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
	storedUser, err := a.DB.GetUserByUsername(user.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Создание JWT-токена
	expirationTime := time.Now().Add(24 * time.Hour) // Срок действия токена - 24 часа
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("secret-key")) // Замените "secret-key" на ваш секретный ключ
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Отправка токена в ответе
	response := map[string]string{
		"token": signedToken,
	}
	json.NewEncoder(w).Encode(response)
}

func (a *API) protectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Проверка наличия токена в заголовке Authorization
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization token required", http.StatusUnauthorized)
		return
	}

	// Проверка и верификация токена
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret-key"), nil // Замените "secret-key" на ваш секретный ключ
	})
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Проверка валидности токена
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Добавьте здесь вашу бизнес-логику для защищенного маршрута
	response := map[string]string{
		"message":  "Protected resource accessed",
		"username": claims.Username,
	}
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

func (a *API) helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := model.Message{Text: "Hello, World!"}
	json.NewEncoder(w).Encode(response)
}

func (a *API) messageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request model.Message
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Сохранение сообщения в базу данных
	err = a.DB.SaveMessage(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(request)
}

func (a *API) RunServer() {
	http.HandleFunc("/", a.helloHandler)
	http.HandleFunc("/message", a.messageHandler)
	http.HandleFunc("/register", a.registerHandler)
	http.HandleFunc("/login", a.loginHandler)
	http.HandleFunc("/protected", a.protectedHandler)
	http.HandleFunc("/logout", a.logoutHandler)

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
