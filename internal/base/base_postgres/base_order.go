package base_postgres

import (
	"application_template/utils"
	"gorm.io/gorm"
	"strings"
)

type OrderFilter interface {
	processValues(db *gorm.DB) []string
	sort() func(db *gorm.DB) *gorm.DB
}

type Order struct {
	Values []string
	Model  HasId
}

func NewOrder(values []string, model HasId) *Order {
	return &Order{
		Values: values,
		Model:  model,
	}
}

func (o *Order) processValues(db *gorm.DB) []string {
	var results []string
	var order string = "asc"

	for _, value := range o.Values {
		attribute := utils.ToSnakeCase(value)

		if strings.HasPrefix(attribute, "-") {
			order = "desc"
			attribute = strings.TrimPrefix(attribute, "-")
		}

		if db.Migrator().HasColumn(o.Model, attribute) {
			results = append(results, strings.Join([]string{attribute, order}, " "))
		}
	}

	return results
}

func (o *Order) sort() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		values := o.processValues(db)

		if len(values) == 0 {
			return db
		}
		return db.Order(strings.Join(values, ", "))
	}
}
