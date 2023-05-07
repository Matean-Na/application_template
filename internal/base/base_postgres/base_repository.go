package base_postgres

import (
	"application_template/internal/database/connect"
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PgError struct {
	Init error
	Pg   error
}

type CrudRepository interface {
	FindAll(p Pager, o OrderFilter, s Scope, total *int64, a interface{}, se Searcher) error
	FindAllDeleted(p Pager, o OrderFilter, s Scope, total *int64, a interface{}, se Searcher) error
	GetFull(s Scope, a interface{}) error
	FindOne(id uint, s Scope, o interface{}) error
	FindOneDeleted(id uint, s Scope, o interface{}) error
	Create(s func(*gorm.DB) *gorm.DB, i interface{}) error
	Update(id uint, s func(*gorm.DB) *gorm.DB, o, u interface{}) error
	Delete(entity HasId) error
	Save(entity HasId) error
	PartialUpdate(entity HasId) error
	Recover(entity HasId) error
	CreateOrUpdate(s func(db *gorm.DB) *gorm.DB, i interface{}, cons string, cols []string) error
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
	FindWhere(o interface{}, w ...interface{}) error
	Transaction(fc func(tx *gorm.DB) error) error

	Raw(sql string, values ...interface{}) (tx *gorm.DB)
	Row() *sql.Row
}

type CrudRepo struct {
	db *gorm.DB
}

func New() CrudRepository {
	return &CrudRepo{
		db: connect.PostgresDB,
	}
}

func (cr *CrudRepo) FindAll(p Pager, o OrderFilter, s Scope, total *int64, a interface{}, se Searcher) error {
	if res := cr.db.Model(a).Scopes(s).Count(total); res.Error != nil {
		return res.Error
	}

	if res := cr.db.Scopes(p.paginate(), o.sort(), s).Find(a); res.Error != nil {
		return res.Error
	}
	if se != nil {
		if res := cr.db.Where(se.getQueryJoin()).Joins(se.getJoinModels()).Scopes(p.paginate(), o.sort(), s).Find(a); res.Error != nil {
			return res.Error
		}
	}
	return nil
}
func (cr *CrudRepo) FindAllDeleted(p Pager, o OrderFilter, s Scope, total *int64, a interface{}, se Searcher) error {
	if res := cr.db.Unscoped().Where("deleted_at IS NOT NULL").Model(a).Scopes(s).Count(total); res.Error != nil {
		return res.Error
	}

	if res := cr.db.Unscoped().Where("deleted_at IS NOT NULL").Scopes(p.paginate(), o.sort(), s).Find(a); res.Error != nil {
		return res.Error
	}
	if se != nil {
		if res := cr.db.Unscoped().Where("deleted_at IS NOT NULL").Where(se.getQueryJoin()).Joins(se.getJoinModels()).Scopes(p.paginate(), o.sort(), s).Find(a); res.Error != nil {
			return res.Error
		}
	}
	return nil
}

func (cr *CrudRepo) GetFull(s Scope, a interface{}) error {
	if res := cr.db.Model(a).Scopes(s); res.Error != nil {
		return res.Error
	}

	if res := cr.db.Scopes(s).Find(a); res.Error != nil {
		return res.Error
	}
	return nil
}

func (cr *CrudRepo) FindOne(id uint, s Scope, o interface{}) error {
	if res := cr.db.Scopes(s).Where("id = ?", id).First(o); res.Error != nil {
		return res.Error
	}
	return nil
}

func (cr *CrudRepo) FindOneDeleted(id uint, s Scope, o interface{}) error {
	if res := cr.db.Unscoped().Where("deleted_at IS NOT NULL").Scopes(s).Where("id = ?", id).First(o); res.Error != nil {
		return res.Error
	}
	return nil
}

func (cr *CrudRepo) Create(s func(*gorm.DB) *gorm.DB, i interface{}) error {
	if res := cr.db.Create(i); res.Error != nil {
		return res.Error
	}
	if res := cr.db.Scopes(s).Model(i).First(i); res.Error != nil {
		return res.Error
	}
	return nil
}

func (cr *CrudRepo) Update(id uint, s func(*gorm.DB) *gorm.DB, o, u interface{}) error {
	if res := cr.db.Where("id = ?", id).First(o); res.Error != nil {
		return res.Error
	}
	if res := cr.db.Model(o).Updates(u); res.Error != nil {
		return res.Error
	}
	if res := cr.db.Scopes(s).Where("id = ?", id).First(o); res.Error != nil {
		return res.Error
	}
	return nil
}

func (cr *CrudRepo) Delete(entity HasId) error {
	if err := cr.FindOne(entity.GetId(), NoScope, entity); err != nil {
		return err
	}
	if res := cr.db.Delete(entity); res.Error != nil {
		return res.Error
	}
	return nil
}

func (cr *CrudRepo) Save(entity HasId) error {
	res := cr.db.Save(entity)
	return res.Error
}

func (cr *CrudRepo) PartialUpdate(entity HasId) error {
	res := cr.db.Updates(entity)
	return res.Error
}
func (cr *CrudRepo) Recover(entity HasId) error {
	if err := cr.FindOneDeleted(entity.GetId(), NoScope, entity); err != nil {
		return err
	}

	if res := cr.db.Unscoped().Where("deleted_at IS NOT NULL").Model(&entity).Update("deleted_at", nil); res.Error != nil {
		return res.Error
	}

	return nil
}

func (cr *CrudRepo) CreateOrUpdate(s func(db *gorm.DB) *gorm.DB, i interface{}, cons string, cols []string) error {
	res := cr.db.Debug().Clauses(clause.OnConflict{
		OnConstraint: cons,
		DoUpdates:    clause.AssignmentColumns(cols),
	}).Create(i)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (cr *CrudRepo) Where(query interface{}, args ...interface{}) (tx *gorm.DB) {
	tx = cr.db
	return tx.Where(query, args)
}

func (cr *CrudRepo) FindWhere(o interface{}, w ...interface{}) error {
	tx := cr.db
	for _, it := range w {
		tx = tx.Where(it)
	}

	if res := tx.Find(o); res.Error != nil {
		return res.Error
	}

	return nil
}

func (cr *CrudRepo) Transaction(fc func(tx *gorm.DB) error) error {
	return cr.db.Transaction(fc)
}

func (cr *CrudRepo) Raw(sql string, values ...interface{}) (tx *gorm.DB) {
	return cr.db.Raw(sql, values)
}

func (cr *CrudRepo) Row() *sql.Row {
	return cr.db.Row()
}
