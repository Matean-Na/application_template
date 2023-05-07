package base_postgres

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	put   = "PUT"
	patch = "PATCH"
)

type ModelInterface interface {
	GetAll() interface{}
	GetOne() HasId

	ScopeAll(db *gorm.DB) *gorm.DB
	ScopeOne(db *gorm.DB) *gorm.DB
}

type CrudInterface interface {
	Register(r *gin.RouterGroup, s string) *gin.RouterGroup

	FindAll(ctx *gin.Context) *AppError
	FindOne(cx *gin.Context) *AppError
	Create(ctx *gin.Context) *AppError
	Update(ctx *gin.Context) *AppError
	Delete(ctx *gin.Context) *AppError
}

type RedisInterface interface {
	KeyAll() string
	KeyOne(id string) string
}

type CrudController struct {
	CrudInterface  CrudInterface
	ModelInterface ModelInterface
	Service        CrudServiceInterface
}

func NewCrudController() *CrudController {
	crudController := &CrudController{}
	crudController.CrudInterface = crudController
	crudController.Service = NewCrudService()
	return crudController
}

func (cc *CrudController) Register(r *gin.RouterGroup, s string) *gin.RouterGroup {
	g := r.Group(s)
	g.GET("", AppHandler(cc.CrudInterface.FindAll).Handle)
	g.GET(":id", AppHandler(cc.CrudInterface.FindOne).Handle)
	g.POST("", AppHandler(cc.CrudInterface.Create).Handle)
	g.PATCH(":id", AppHandler(cc.CrudInterface.Update).Handle)
	g.DELETE(":id", AppHandler(cc.CrudInterface.Delete).Handle)
	return g
}

func (cc *CrudController) FindAll(ctx *gin.Context) *AppError {
	ct := NewCrudTemplate(cc)
	return ct.FindAll(ctx, cc.Service)
}

func (cc *CrudController) FindOne(ctx *gin.Context) *AppError {
	ct := NewCrudTemplate(cc)
	return ct.FindOne(ctx, cc.Service)
}

func (cc *CrudController) Create(ctx *gin.Context) *AppError {
	ct := NewCrudTemplate(cc)
	return ct.Create(ctx, cc.Service)
}

func (cc *CrudController) Update(ctx *gin.Context) *AppError {
	var partial bool
	ct := NewCrudTemplate(cc)
	switch ctx.Request.Method {
	case put:
		partial = false
	case patch:
		partial = true
	}
	return ct.Update(ctx, cc.Service, partial)
}

func (cc *CrudController) Delete(ctx *gin.Context) *AppError {
	ct := NewCrudTemplate(cc)
	return ct.Delete(ctx, cc.Service)
}

func (cc *CrudController) KeyAll() string {
	return fmt.Sprintf("%T:all", cc.CrudInterface)
}

func (cc *CrudController) KeyOne(id string) string {
	return fmt.Sprintf("%T:%s", cc.CrudInterface, id)
}
