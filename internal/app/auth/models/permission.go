package models

import "application_template/internal/base/base_postgres"

type Permission struct {
	base_postgres.Entity
	IdRole uint   `gorm:"index:idx_permission_unique,unique,where:deleted_at is null"`
	Role   Role   `gorm:"foreignKey:IdRole"`
	Type   uint   `gorm:"index:idx_permission_unique"`
	Target string `gorm:"index:idx_permission_unique"`
	Value  uint
}
