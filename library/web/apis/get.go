package apis

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"

	"net/http"
	"startkit/library/gins"
	"startkit/library/gorms"
	"startkit/library/random"
	"startkit/library/structs"

	"github.com/gin-gonic/gin"
)

type GET struct {
	API
	Conditions  map[string]string
	Comparisons map[string]map[string]string
	DBWheres    []func(db *gorm.DB) *gorm.DB
}

type GETFunc func(g *GET) func(db *gorm.DB) *gorm.DB

// GETComparisonFunc conditions contain the db column as the 1st key,
// the comparison string as the 2nd key, the value in the map is the value for query
type GETComparisonFunc func(g *GET, key string, conditions map[string]string) func(db *gorm.DB) *gorm.DB

var (
	GETGeneralFuncs = map[string]GETFunc{
		"limit":        Pagination,
		"start":        Pagination,
		"order_by":     OrderBy,
		"is_random_id": RandomID,
	}
	// GETComparisonFuncs default use AND in sql statment for multiple statments for comparison
	GETComparisonFuncs = map[string]GETComparisonFunc{
		"gt":   GreaterThan,
		"st":   SmallerThan,
		"gtoe": GreaterThanOrEqual,
		"stoe": SmallerThanOrEqual,
		"like": Like,
	}
)

func (g *GET) Run() (err error) {
	if g.StaticPath != "" {
		g.ServeStatic()
		return
	}
	// g.DBResult = reflect.Zero(reflect.ValueOf(g.DBResult).Type()).Elem().Interface() // mistake
	// new a pointer Value p of the DBResult (pointer of the structure that use for query)
	p := reflect.ValueOf(g.DBResult).Elem()
	// set the pointer
	p.Set(reflect.Zero(p.Type()))
	g.FillInIDer()
	g.FillInConditions()
	g.FillInComparisons()
	for _, name := range g.API.Context.App.InUseService {
		switch name {
		case "Mysql":
			err = g.MysqlHandler()
		}
	}
	return
}

func (g *GET) ServeStatic() {
	g.Ctx.HTML(
		http.StatusOK,
		g.StaticPath,
		gin.H{
			"response": Resp("Success", "Serve Static HTML"),
		},
	)
}

func (g *GET) FillInIDer() {
	if g.IDer != "" {
		if id := g.Ctx.Param(g.IDer); id != "" {
			g.DBWheres = append(g.DBWheres, gorms.Where(g.IDer, " = ? ", id))
		}
	}
	return
}

func (g *GET) FillInConditions() {
	g.Conditions = make(map[string]string)
	jsons := structs.GetTags(g.Structure, "json")
	gormColumns := structs.GetTagsValueWithSpliter(g.Structure, "gorm", ":")
	if values := gins.URLParam(g.Params, g.Ctx); len(values) > 0 {
		for key, value := range values {
			if value != "" {
				g.Conditions[key] = value
			}
		}
	}
	if values := gins.URLDefaultQuery(g.Querys, g.Ctx); len(values) > 0 {
		for key, value := range values {
			if value != "" {
				g.Conditions[key] = value
			}
		}
	}
	for k, v := range g.Conditions {
		if _, ok := jsons[k]; ok {
			if _, ok := gormColumns[k]; ok {
				g.DBWheres = append(g.DBWheres, gorms.Where(k, " IN (?) ", v))
			}
		}
		if _, ok := GETGeneralFuncs[k]; ok {
			if generalFunc := GETGeneralFuncs[k](g); generalFunc != nil {
				g.DBWheres = append(g.DBWheres, generalFunc)
			}
		}
	}
	return
}

func (g *GET) FillInComparisons() {
	g.Comparisons = make(map[string]map[string]string)
	jsons := structs.GetTags(g.Structure, "json")
	gormColumns := structs.GetTagsValueWithSpliter(g.Structure, "gorm", ":")
	if values := gins.URLQueryMap(g.Querys, g.Ctx); len(values) > 0 {
		for key, value := range values {
			if len(value) > 0 {
				g.Comparisons[key] = make(map[string]string)
				g.Comparisons[key] = value
			}
		}
	}
	for k, v := range g.Comparisons {
		if _, ok := jsons[k]; ok {
			if _, ok := gormColumns[k]; ok {
				for statement, value := range v {
					if len(value) > 0 {
						if comparisonFunc := GETComparisonFuncs[statement](g, k, map[string]string{statement: value}); comparisonFunc != nil {
							g.DBWheres = append(g.DBWheres, comparisonFunc)
						}
					}
				}
			}
		}
	}
	return
}

func (g *GET) MysqlHandler() error {
	if count := len(g.API.ValidatorFuncs); count > 0 {
		for i := 0; i < count; i++ {
			if code, message, err := g.API.ValidatorFuncs[i](&g.API); code > 0 || message != "" {
				g.Ctx.AbortWithStatusJSON(code, gin.H{"response": Resp(message, err.Error())})
				return err
			}
		}
	}
	if count, err := gorms.TotalCount(&g.Mysql, g.DBObject, nil); err != nil && err != gorm.ErrRecordNotFound {
		g.Ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"response": Resp("DB Error", err.Error())})
		return err
	} else if count <= 0 || err == gorm.ErrRecordNotFound {
		g.Ctx.JSON(http.StatusOK, gin.H{"response": Resp("No Record", map[string]int{"count": count})})
		return nil
	}
	if err := gorms.ScopesQuery(&g.Mysql, g.DBWheres, g.DBResult); err != nil {
		g.Ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"response": Resp("DB Error", err.Error())})
		return err
	}
	g.Ctx.JSON(http.StatusOK, gin.H{"response": Resp("Success", g.DBResult)})
	return nil
}

func Pagination(g *GET) func(db *gorm.DB) *gorm.DB {
	if count, err := gorms.TotalCount(&g.Mysql, g.DBObject, nil); err != nil {
		return nil
	} else if count > 0 {
		if _, ok := g.Conditions["start"]; !ok || g.Conditions["start"] == "" {
			g.Conditions["start"] = "0"
		}
		if _, ok := g.Conditions["limit"]; !ok || g.Conditions["limit"] == "" {
			g.Conditions["limit"] = strconv.Itoa(count)
		}
		return gorms.Pagination(g.Conditions["start"], g.Conditions["limit"])
	}
	return nil
}

func RandomID(g *GET) func(db *gorm.DB) *gorm.DB {
	if g.Conditions["is_random_id"] == "true" {
		if count, err := gorms.TotalCount(&g.Mysql, g.DBObject, nil); err != nil {
			return nil
		} else if count > 0 {
			return gorms.Where("id", " IN (?) ", random.RandIntArray(0, count, count))
		}
	}
	return nil
}

func OrderBy(g *GET) func(db *gorm.DB) *gorm.DB {
	if count, err := gorms.TotalCount(&g.Mysql, g.DBObject, nil); err != nil {
		return nil
	} else if count > 0 {
		if _, ok := g.Conditions["order_by"]; !ok || g.Conditions["order_by"] == "" {
			g.Conditions["order_by"] = "id"
		}
	}
	return nil
}

func GreaterThan(g *GET, key string, conditions map[string]string) func(db *gorm.DB) *gorm.DB {
	return gorms.Where(key, " > ?", conditions["gt"])
}

func SmallerThan(g *GET, key string, conditions map[string]string) func(db *gorm.DB) *gorm.DB {
	return gorms.Where(key, " < ?", conditions["st"])
}

func GreaterThanOrEqual(g *GET, key string, conditions map[string]string) func(db *gorm.DB) *gorm.DB {
	return gorms.Where(key, " >= ?", conditions["gtoe"])
}

func SmallerThanOrEqual(g *GET, key string, conditions map[string]string) func(db *gorm.DB) *gorm.DB {
	return gorms.Where(key, " <= ?", conditions["stoe"])
}

func Like(g *GET, key string, conditions map[string]string) func(db *gorm.DB) *gorm.DB {
	return gorms.Where(key, " LIKE ?", "%"+strings.Trim(conditions["like"], "\"")+"%")
}
