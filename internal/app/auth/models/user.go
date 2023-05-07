package models

import (
	"application_template/internal/base/base_postgres"
	"errors"
	"html"
	"strings"
	"time"
)

type User struct {
	base_postgres.Entity
	UserName     string `gorm:"index:idx_user_unique,unique,where:deleted_at is null"`
	UserPassword string
	Active       bool
	Roles        []Role `gorm:"many2many:user_roles;"`
}

//func (u *User) BeforeCreate(*gorm.DB) error {
//	hashedPassword, err := security.Hash(u.UserPassword)
//	if err != nil {
//		return err
//	}
//	u.UserPassword = string(hashedPassword)
//	return nil
//}

func (u *User) Prepare() {
	u.UserName = html.EscapeString(strings.TrimSpace(u.UserName))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "login":
		if u.UserName == "" {
			err = errors.New("required Username")
			errorMessages["Required_username"] = err.Error()
		}
		if u.UserPassword == "" {
			err = errors.New("required Password")
			errorMessages["Required_password"] = err.Error()
		}
	}
	return errorMessages
}
