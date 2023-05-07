package handlers

import (
	"application_template/internal/app/auth/models"
	"application_template/internal/app/auth/services"
	"encoding/json"
	"log"
	"net/http"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.User

	// Разбор данных запроса
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Вызов сервиса для регистрации нового пользователя
	err = h.authService.Registration(&dto)
	if err != nil {
		log.Println("Error registering user:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Отправка успешного ответа
	w.WriteHeader(http.StatusOK)
}
