package base_postgres

import (
	"application_template/internal/database/connect"
	"application_template/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type AppError struct {
	Error     string
	Code      int
	Message   string
	Detailed  string
	FieldName string
	Fields    map[string]string
}

type AppHandler func(ctx *gin.Context) *AppError

func (a AppHandler) Handle(ctx *gin.Context) {
	if err := a(ctx); err != nil {
		ctx.JSON(err.Code, gin.H{"error": err})
	}
}

func Ok(ctx *gin.Context, i interface{}) *AppError {
	ctx.JSON(http.StatusOK, i)
	return nil
}

func ErrNotExist(ctx *gin.Context, i interface{}) *AppError {
	ctx.JSON(http.StatusNotFound, i)
	return nil
}

func OkT(ctx *gin.Context, t int64, i interface{}) *AppError {
	ctx.JSON(http.StatusOK, gin.H{"total": t, "data": i})
	return nil
}

func okJson(ctx *gin.Context, s string) *AppError {
	ctx.Data(http.StatusOK, "application/json; charset=utf-8", []byte(s))
	return nil
}

func okMessage(ctx *gin.Context, SuccessMessage string, Message string) *AppError {
	ctx.JSON(http.StatusOK, gin.H{"Success": SuccessMessage, "Message": Message})
	return nil
}

func errBind(err error) *AppError {
	return &AppError{
		Error:     err.Error(),
		Code:      http.StatusBadRequest,
		Message:   "Bind error",
		Detailed:  "",
		FieldName: "",
	}
}

func ErrNotFound(ctx *gin.Context, err error, instance interface{}, id int, params string) *AppError {
	var message string
	if params != "" {
		message = utils.Localize(
			ctx,
			"exception:failed-to-fetch-records-with-params",
			map[string]interface{}{
				"Table":      GetTableName(instance, connect.PostgresDB),
				"Parameters": params,
			},
		)
	} else {
		if id != 0 {
			message = utils.Localize(ctx, "exception:failed-to-fetch-one-record", map[string]interface{}{
				"Table": GetTableName(instance, connect.PostgresDB),
				"ID":    id,
			})
		} else {
			message = utils.Localize(ctx, "exception:could-not-fetch-records", map[string]string{
				"Table": GetTableName(instance, connect.PostgresDB),
			})
		}
	}

	return &AppError{
		Error:     err.Error(),
		Code:      http.StatusNotFound,
		Message:   message,
		Detailed:  "",
		FieldName: "",
	}
}

func ErrNotCreated(ctx *gin.Context, err error, instance interface{}) *AppError {
	message := utils.Localize(ctx, "exception:failed-to-create-record", map[string]interface{}{
		"Table": GetTableName(instance, connect.PostgresDB),
	})
	return &AppError{
		Error:     err.Error(),
		Code:      http.StatusBadRequest,
		Message:   message,
		Detailed:  "",
		FieldName: "",
	}
}

func ErrNotUpdated(err error) *AppError {
	return &AppError{
		Error:     err.Error(),
		Code:      http.StatusBadRequest,
		Message:   "Record not updated",
		Detailed:  "",
		FieldName: "",
	}
}

func errNotDeleted(err error) *AppError {
	return &AppError{
		Error:     err.Error(),
		Code:      http.StatusBadRequest,
		Message:   "Record not deleted",
		Detailed:  "",
		FieldName: "",
	}
}

func ErrDenied(err error) *AppError {
	return &AppError{
		Error:     err.Error(),
		Code:      http.StatusForbidden,
		Message:   "Access denied",
		Detailed:  "",
		FieldName: "",
	}
}

func ErrBadRequest(ctx *gin.Context, err error, data map[string]interface{}) *AppError {
	return &AppError{
		Error:     err.Error(),
		Code:      http.StatusBadRequest,
		Message:   utils.Localize(ctx, err.Error(), data),
		Detailed:  "",
		FieldName: "",
	}
}

func LocalizeError(ctx *gin.Context, err error) *AppError {
	var pgxError *pgconn.PgError
	var fieldName string
	var detailedInfo string
	switch e := err.(type) {
	case utils.LocalizeError:
		errors.As(e.Source, &pgxError)
		if pgxError != nil {
			fieldName = ParseFieldName(pgxError.Detail)
			detailedInfo = pgxError.Detail
		}
		if pgxError != nil {
			fieldName = ParseFieldName(pgxError.Detail)
			detailedInfo = pgxError.Detail
		}
		return &AppError{
			Error:     e.Error(),
			Message:   e.Localize(ctx),
			Code:      http.StatusBadRequest,
			Detailed:  detailedInfo,
			FieldName: fieldName,
		}
	case *utils.LocalizeError:
		errors.As(e.Source, &pgxError)
		if pgxError != nil {
			fieldName = ParseFieldName(pgxError.Detail)
			detailedInfo = pgxError.Detail
		}
		return &AppError{
			Error:     e.Error(),
			Message:   e.Localize(ctx),
			Code:      http.StatusBadRequest,
			Detailed:  detailedInfo,
			FieldName: fieldName,
		}
	}
	errors.As(err, &pgxError)
	if pgxError != nil {
		fieldName = ParseFieldName(pgxError.Detail)
		detailedInfo = pgxError.Detail
	}
	return DefaultError(ctx, err, fieldName, detailedInfo)
}

func ParseFieldName(input string) string {
	i := strings.Index(input, "(")
	if i >= 0 {
		j := strings.Index(input, ")")
		if j >= 0 {
			return toCamelCase(input[i+1 : j])
		}
	}
	return ""
}

func toCamelCase(input string) string {
	words := strings.Split(input, "_")
	key := strings.Title(words[0])
	for _, word := range words[1:] {
		key += strings.Title(word)
	}
	return key
}

func DefaultError(ctx *gin.Context, err error, fieldName string, detailedInfo string) *AppError {
	return &AppError{
		Error:     err.Error(),
		Message:   utils.DefaultError(ctx),
		Code:      http.StatusBadRequest,
		Detailed:  detailedInfo,
		FieldName: fieldName,
	}
}

func ParamInt(p string) int {
	val, err := strconv.Atoi(p)
	if err != nil {
		return 0
	}
	return val
}

func ParamDecimal(p string) decimal.Decimal {
	val, err := strconv.ParseFloat(p, 64)
	if err != nil {
		return decimal.NewFromInt(0)
	}
	return decimal.NewFromFloat(val)
}

func ParamUint(p string) uint {
	return uint(ParamInt(p))
}

func getPager(ctx *gin.Context) Pager {
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return NewPager(page, pageSize, offset)
}

func CheckDate(format, date string) bool {
	t, err := time.Parse(format, date)
	if err != nil {
		return false
	}
	return t.Format(format) == date
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func GetBody(ctx *gin.Context, o interface{}) error {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, o)
	if err != nil {
		return err
	}

	return nil
}

func getOrder(ctx *gin.Context, model HasId) OrderFilter {
	var values []string = ctx.Request.URL.Query()["order_by"]
	return NewOrder(values, model)
}

func getQuery(c *gin.Context, a interface{}) Searcher {
	search := c.Query("search")
	if search != "" {
		in := []byte(search)
		var raw map[string]interface{}
		if err := json.Unmarshal(in, &raw); err != nil {
			return nil
		}

		data := reflect.ValueOf(a)
		s := data.Type()

		model := fmt.Sprintf("%v", s)
		myModel := model[10:] + "s"

		var query string
		var queryJoin string
		var isJoin bool
		var JoinModels string
		i := 1
		lenRaw := len(raw)
		for key, v := range raw {
			value := fmt.Sprintf("%v", v)
			typeValue := fmt.Sprintf("%T", v)
			var fromV string
			var beforeV string
			var lang string
			var word string

			if strings.Contains(key, "json") == true {
				s := strings.Split(value, " = ")
				lang = s[0]
				word = s[1]
				typeValue = "json"
			}

			if CheckDate("2006-01-02", value) == true {
				typeValue = "date"
			}

			if strings.Contains(value, "and") == true {
				s := strings.Split(value, " and ")
				fromV = s[0]
				beforeV = s[1]
				if CheckDate("2006-01-02", fromV) == true {
					typeValue = "date"
				}
			}

			if typeValue == "map[string]interface {}" {
				byteV, _ := json.Marshal(v)
				var rawV map[string]interface{}
				if err := json.Unmarshal(byteV, &rawV); err != nil {
					return nil
				}
				JoinModels = key
				for k, val := range rawV {
					newVal := fmt.Sprintf("%v", val)
					typeVal := fmt.Sprintf("%T", val)
					var from string
					var before string

					if CheckDate("2006-02-01", newVal) == true {
						typeVal = "date"
					}

					if strings.Contains(newVal, "and") == true {
						s := strings.Split(newVal, " and ")
						from = s[0]
						before = s[1]
						if CheckDate("2006-02-01", from) == true {
							typeVal = "date"
						}
					}

					if i < lenRaw {
						if typeVal == "bool" {
							queryJoin += k + " = " + newVal + " and "
						} else if typeVal == "string" && newVal != "not_null" {
							queryJoin += k + " ILIKE " + "'%" + newVal + "%' and "
						}
						if typeVal == "int" || typeVal == "float64" {
							queryJoin += k + " = " + newVal + " and "
						}
						if typeVal == "date" {
							if strings.Contains(newVal, "and") == true {
								queryJoin += k + " BETWEEN '" + from + "' AND '" + before + "' and "
							} else {
								queryJoin += k + "::text LIKE " + "'%" + newVal + "%' and "
							}
						}
						if newVal == "not_null" {
							queryJoin += k + " IS NOT NULL and "
						}
					}
					if i >= lenRaw {
						if typeVal == "bool" {
							queryJoin += k + " = " + newVal
						} else if typeVal == "string" && newVal != "not_null" {
							queryJoin += k + " ILIKE " + "'%" + newVal + "%'"
						}
						if typeVal == "int" || typeVal == "float64" {
							queryJoin += k + " = " + newVal
						}
						if typeVal == "date" {
							if strings.Contains(newVal, "and") == true {
								queryJoin += k + " BETWEEN '" + from + "' AND '" + before + "'"
							} else {
								queryJoin += k + "::text LIKE " + "'%" + newVal + "%'"
							}
						}
						if newVal == "not_null" {
							queryJoin += k + " IS NOT NULL"
						}
					}
					i++
				}
				isJoin = true
			} else {
				if i < lenRaw {
					if typeValue == "bool" {
						queryJoin += ToSnakeCase(myModel) + "." + key + " = " + value + " and "
					} else if typeValue == "string" && value != "not_null" {
						queryJoin += ToSnakeCase(myModel) + "." + key + " ILIKE " + "'%" + value + "%' and "
					}
					if typeValue == "int" || typeValue == "float64" {
						queryJoin += ToSnakeCase(myModel) + "." + key + " = " + value + " and "
					}
					if typeValue == "date" {
						if strings.Contains(value, "and") == true {
							queryJoin += key + " BETWEEN '" + fromV + "' AND '" + beforeV + "' and "
						} else {
							queryJoin += ToSnakeCase(myModel) + "." + key + "::text LIKE " + "'%" + value + "%' and "
						}
					}
					if value == "not_null" {
						queryJoin += ToSnakeCase(myModel) + "." + key + " IS NOT NULL and "
					}
					if typeValue == "json" {
						queryJoin += ToSnakeCase(myModel) + "." + key + " ->> '" + lang + "' ilike '%" + word + "%' and "
					}
				}
				if i >= lenRaw {
					if typeValue == "bool" {
						queryJoin += ToSnakeCase(myModel) + "." + key + " = " + value
					} else if typeValue == "string" && value != "not_null" {
						queryJoin += ToSnakeCase(myModel) + "." + key + " ILIKE " + "'%" + value + "%'"
					}
					if typeValue == "int" || typeValue == "float64" {
						queryJoin += ToSnakeCase(myModel) + "." + key + " = " + value
					}
					if typeValue == "date" {
						if strings.Contains(value, "and") == true {
							queryJoin += key + " BETWEEN '" + fromV + "' AND '" + beforeV + "'"
						} else {
							queryJoin += ToSnakeCase(myModel) + "." + key + "::text LIKE " + "'%" + value + "%'"
						}
					}
					if value == "not_null" {
						queryJoin += ToSnakeCase(myModel) + "." + key + " IS NOT NULL"
					}
					if typeValue == "json" {
						queryJoin += ToSnakeCase(myModel) + "." + key + " ->> '" + lang + "' ilike '%" + word + "%'"
					}
				}
				i++
			}
		}
		return NewSearcher(query, isJoin, JoinModels, queryJoin)
	} else {
		return nil
	}
}

func I18nError(c *gin.Context, model interface{}, errCode string) *AppError {
	codes := map[string]int{
		"exception:could-not-count-records": http.StatusBadRequest,
		"exception:could-not-fetch-records": http.StatusNotFound,
		"exception:failed-to-create-record": http.StatusBadRequest,
		"exception:failed-to-update-record": http.StatusBadRequest,
		"exception:failed-to-delete-record": http.StatusBadRequest,
	}

	var code = http.StatusBadRequest
	if c, found := codes[errCode]; found {
		code = c
	}

	table := GetTableName(model, connect.PostgresDB)
	return &AppError{
		Error: errCode,
		Message: utils.Localize(c, errCode, map[string]interface{}{
			"Table": table,
		}),
		Code: code,
	}
}
