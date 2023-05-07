package services

import (
	"application_template/internal/app/auth/models"
	"application_template/internal/app/auth/repositories"
)

type AuthService struct {
	authRepo repositories.AuthRepository
}

func NewAuthService(authRepo repositories.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

// Registration регистрирует нового пользователя.
func (s *AuthService) Registration(user *models.User) error {
	// Валидация данных
	// бизнес логика
	return s.authRepo.Registration(user)
}
