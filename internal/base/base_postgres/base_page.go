package base_postgres

import "gorm.io/gorm"

type Pager interface {
	getPage() int
	getPageSize() int
	getOffset() int

	paginate() func(db *gorm.DB) *gorm.DB
}

type Page struct {
	page, pageSize, offset int
}

func NewPager(page, pageSize, offset int) Pager {
	return &Page{
		page:     page,
		pageSize: pageSize,
		offset:   offset,
	}
}

func (p *Page) getPage() int {
	return p.page
}

func (p *Page) getPageSize() int {
	return p.pageSize
}

func (p *Page) getOffset() int {
	return p.offset
}

func (p *Page) paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.getOffset()).Limit(p.getPageSize())
	}
}
