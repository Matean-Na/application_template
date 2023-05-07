package models

import (
	"application_template/internal/base/base_postgres"
)

type Role struct {
	base_postgres.Entity
	Name  string `gorm:"index:idx_role_unique,unique,where:deleted_at is null"`
	Users []User `gorm:"many2many:user_roles;"`
}

//func (t *Role) BeforeCreate(tx *gorm.DB) error {
//	var role []Role
//	if err := connect.DB.Raw("select * from roles;").Scan(&role).Error; err != nil {
//		return err
//	}
//
//	for _, v := range role {
//		var exist bool
//		if t.Name == v.Name {
//			exist = true
//		}
//		if v.DeletedAt.Valid == true && exist == true {
//			if err := connect.DB.Raw("update roles set deleted_at = null where id = ?", v.ID).Scan(&role).Error; err != nil {
//				return err
//			}
//			return errors.New("exception: deleted record restored")
//		} else if exist == true && v.DeletedAt.Valid == false {
//			return errors.New("exception:record-already-exist")
//		}
//	}
//	return nil
//}
