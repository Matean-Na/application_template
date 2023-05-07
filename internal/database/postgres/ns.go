package postgres

import (
	"fmt"
	"gorm.io/gorm/schema"
	"strings"
)

type ProNamingStrategy struct {
	schema.NamingStrategy
}

func (ns ProNamingStrategy) RelationshipFKName(rel schema.Relationship) string {
	return ns.formatName("fk", rel.Schema.Table, ns.ColumnName("tab", rel.Name))
}

func (ns ProNamingStrategy) IndexName(table, column string) string {
	return ns.formatName("idx", table, ns.ColumnName(table, column))
}

func (ns ProNamingStrategy) formatName(prefix, table, name string) string {
	formattedName := strings.Replace(fmt.Sprintf("%v_%v_%v", prefix, table, name), ".", "_", -1)
	return formattedName
}
