package postgres

import (
	"application_template/internal/app/auth/models"
	"application_template/internal/base/base_postgres"
	"application_template/internal/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDsn(config config.DB) string {
	return fmt.Sprintf("host=%s user=%s password=%s DBName=%s port=%s sslmode=%s TimeZone=UTC",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort, config.DBSSLMode)
}

func Recreate(config config.DB, sysDb string) error {
	sysConf := config
	sysConf.DBName = sysDb
	dsn := GetDsn(sysConf)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to %s database", sysConf.DBName)
	}

	res := db.Exec(fmt.Sprintf("drop database if exists %s", config.DBName))
	if res.Error != nil {
		return fmt.Errorf("failed to drop %s database", config.DBName)
	}

	res = db.Exec(fmt.Sprintf("create database %s", config.DBName))
	if res.Error != nil {
		return fmt.Errorf("failed to create %s database", config.DBName)
	}

	sqlDb, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sqldb")
	}

	if err = sqlDb.Close(); err != nil {
		return fmt.Errorf("failed to close sqldb")
	}

	return nil
}

var Models = []interface{}{}

func notAll(db *gorm.DB) bool {
	res := true

	for _, model := range Models {
		e := db.Migrator().HasTable(model)
		res = res && !e
	}

	return res
}

func createConstraint(database *gorm.DB, cons string, m interface{}, cols string) error {
	t := base_postgres.GetTableName(m, database)
	exist := database.Migrator().HasConstraint(m, cons)
	if !exist {
		sql := fmt.Sprintf("alter table %s add constraint %s unique (%s)", t, cons, cols)
		if res := database.Exec(sql); res.Error != nil {
			return fmt.Errorf("failed to create constraint error: %s", res.Error)
		}
	}
	return nil
}

func checkHasDuplicateConstraint(database *gorm.DB, index [4]string, m interface{}) error { //
	for _, v := range index {
		existHasDuplicateConstraint := database.Migrator().HasConstraint(m, v)
		if existHasDuplicateConstraint {
			if res := database.Migrator().DropConstraint(m, v); res != nil {
				return fmt.Errorf("failed to delete constraint error: %s", res.Error)
			}
		}
	}
	return nil
}

func Connect(config config.DB) (*gorm.DB, error) {

	dsn := GetDsn(config)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: ProNamingStrategy{},
		//Logger:         logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s database", config.DBName)
	}

	initial := notAll(database)

	duplicateConstraint := [4]string{
		"idx_permissions_id_role",
		"idx_permissions_target",
		"idx_permissions_type",
		"cons_uniq",
	}

	err = checkHasDuplicateConstraint(database, duplicateConstraint, &models.Permission{})
	if err != nil {
		return nil, err
	}

	err = database.AutoMigrate(Models...)
	if err != nil {
		return nil, fmt.Errorf("failed auto migration error: %s", err.Error())
	}

	err = checkHasDuplicateConstraint(database, duplicateConstraint, &models.Permission{})
	if err != nil {
		return nil, err
	}

	if err := database.SetupJoinTable(&models.User{}, "Roles", &models.UserRole{}); err != nil {
		return nil, err
	}

	err = createConstraint(database, "cons_uniq", &models.Permission{}, "id_role,type,target")
	if err != nil {
		return nil, err
	}

	if initial {
		if err = Initialize(database); err != nil {
			return nil, err
		}
	}

	return database, nil
}

func Initialize(db *gorm.DB) error {
	fmt.Println("Initialize Empty Database, loading necessary seeds ...")
	RunInitialDbLoader(db)
	return nil
}

func RunInitialDbLoader(DB *gorm.DB) {
	DB = DB
}
