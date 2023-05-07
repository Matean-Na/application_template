package base_postgres

import (
	"application_template/internal/database/connect"
	"application_template/internal/database/redis"
	"application_template/utils"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"reflect"
)

type CrudTemplateInterface interface {
	FindAll(c *gin.Context) *AppError
	FindOne(c *gin.Context) *AppError
	Create(c *gin.Context) *AppError
	Update(c *gin.Context) *AppError
	Delete(c *gin.Context) *AppError
	GetExcel(c *gin.Context) *AppError
}

type CrudTemplate struct {
	mi ModelInterface
	ri RedisInterface
}

func NewCrudTemplate(cc *CrudController) *CrudTemplate {
	ct := &CrudTemplate{
		mi: cc.ModelInterface,
		ri: cc,
	}
	return ct
}

func (ct *CrudTemplate) FindAllFunc(c *gin.Context, findAll FindAll) *AppError {
	a := ct.mi.GetAll()

	//search := c.Query("search")
	//if search == "" {
	//	all, err := redis.Get(c, ct.ri.KeyAll())
	//	if err == nil {
	//		return okJson(c, all)
	//	}
	//}

	var total int64
	if err := findAll(a, ct.mi.ScopeAll, getPager(c), getOrder(c, ct.mi.GetOne()), &total, getQuery(c, a)); err != nil {
		return I18nError(c, a, "exception:could-not-fetch-records")
	}

	_ = redis.Set(c, ct.ri.KeyAll(), a)

	return OkT(c, total, a)
}

func (ct *CrudTemplate) FindAll(c *gin.Context, allInter FindAllInterface) *AppError {
	return ct.FindAllFunc(c, allInter.FindAll)
}

func (ct *CrudTemplate) FindOneFunc(c *gin.Context, findOne FindOne) *AppError {
	id := c.Param("id")
	o := ct.mi.GetOne()
	o.SetId(ParamUint(id))

	redisStop := c.Query("redisStop")

	if redisStop == "" {
		one, err := redis.Get(c, ct.ri.KeyOne(id))
		if err == nil {
			return okJson(c, one)
		}
	}

	if err := findOne(o, ct.mi.ScopeOne); err != nil {
		return &AppError{
			Error: err.Error(),
			Message: utils.Localize(c, "exception:failed-to-fetch-one-record", map[string]interface{}{
				"Table": GetTableName(o, connect.PostgresDB),
				"ID":    id,
			}),
			Code: http.StatusNotFound,
		}
	}

	_ = redis.Set(c, ct.ri.KeyOne(id), o)

	return Ok(c, o)
}

func (ct *CrudTemplate) FindOne(c *gin.Context, oneInter FindOneInterface) *AppError {
	return ct.FindOneFunc(c, oneInter.FindOne)
}

func (ct *CrudTemplate) CreateFunc(c *gin.Context, create Create) *AppError {
	i := ct.mi.GetOne()

	if err := c.ShouldBindJSON(i); err != nil {
		return LocalizeError(c, err)
	}

	sType := reflect.ValueOf(i).Elem()
	field := sType.FieldByName("IdLanguage")
	if field.IsValid() {
		field.SetUint(uint64(c.GetUint("language")))
	}

	if err := create(i); err != nil {
		return LocalizeError(c, err)
	}

	_ = redis.Unset(c, ct.ri.KeyAll())

	return Ok(c, i)
}

func (ct *CrudTemplate) Create(c *gin.Context, creInter CreateInterface) *AppError {
	return ct.CreateFunc(c, creInter.Create)
}

func (ct *CrudTemplate) UpdateFunc(c *gin.Context, update Update) *AppError {
	id := c.Param("id")
	o := ct.mi.GetOne()

	body, err := ValidateBody(o, c.Request.Body)
	if err != nil {
		return LocalizeError(c, err)
	}

	data, err := json.Marshal(&body)
	if err != nil {
		return LocalizeError(c, utils.NewLocalizeError(
			nil, "exception:marshalling-error", nil,
		))
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(data))

	if err := c.ShouldBindJSON(o); err != nil {
		return LocalizeError(c, err)
	}

	o.SetId(ParamUint(id))

	if err := update(o); err != nil {
		return ErrNotUpdated(err)
	}

	_ = redis.Unset(c, ct.ri.KeyAll())
	_ = redis.Unset(c, ct.ri.KeyOne(id))

	return Ok(c, body)
}

func (ct *CrudTemplate) Update(c *gin.Context, updInter UpdateInterface, partial bool) *AppError {
	if partial {
		return ct.UpdateFunc(c, updInter.PartialUpdate)
	}
	return ct.UpdateFunc(c, updInter.Update)
}

func (ct *CrudTemplate) DeleteFunc(c *gin.Context, delete Delete) *AppError {
	id := c.Param("id")

	o := ct.mi.GetOne()
	o.SetId(ParamUint(id))

	if err := delete(o); err != nil {
		return errNotDeleted(err)
	}

	_ = redis.Unset(c, ct.ri.KeyAll())
	_ = redis.Unset(c, ct.ri.KeyOne(id))

	return Ok(c, o)
}

func (ct *CrudTemplate) Delete(c *gin.Context, delInter DeleteInterface) *AppError {
	return ct.DeleteFunc(c, delInter.Delete)
}
