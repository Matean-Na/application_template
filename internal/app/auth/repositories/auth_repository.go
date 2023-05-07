package repositories

import (
	"application_template/internal/app/auth/dto"
	"application_template/internal/app/auth/models"
	"gorm.io/gorm"
)

type AuthRepositoryInterface interface {
	Registration(user *models.User) error
	Authorisation(auth dto.AuthorizationDTO) (access, refresh string, err error)
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Registration(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *AuthRepository) Authorisation(auth dto.AuthorizationDTO) (access, refresh string, err error) {
	//var user models.User
	//err = r.db.Where("email =?", auth.Email).First(&user).Error
	//if err != nil {
	//	return "", "", err
	//}
	//
	//access, refresh, err = r.db.Model(&user).Association("Access").Get(), r.db.Model(&user).Association("Refresh").Get(), nil
	//return access, refresh, err
	return "", "", nil
}
