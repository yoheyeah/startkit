package apis

import (
	"net/http"
	"startkit/library/gins"
	"startkit/library/gorms"

	"startkit/library/structs"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type DELETE struct {
	API
	Conditions map[string]string
	DBWheres   []func(db *gorm.DB) *gorm.DB
}

func (d *DELETE) Run() (err error) {
	d.FillInIDer()
	d.FillInConditions()
	for _, name := range d.API.Context.App.InUseService {
		switch name {
		case "Mysql":
			err = d.MysqlHandler()
		}
	}
	return
}

func (d *DELETE) FillInIDer() {
	if d.IDer != "" {
		if id := d.Ctx.Param(d.IDer); id != "" {
			if v, ok := d.IDerPropertyRelation[d.IDer]; ok {
				d.DBWheres = append(d.DBWheres, gorms.Where(v, " = ? ", id))
			} else {
				d.DBWheres = append(d.DBWheres, gorms.Where(d.IDer, " = ? ", id))
			}
		}
	}
	return
}

func (d *DELETE) FillInConditions() {
	d.Conditions = make(map[string]string)
	jsons := structs.GetTags(d.Structure, "json")
	gormColumns := structs.GetTagsValueWithSpliter(d.Structure, "gorm", ":")
	if values := gins.URLParam(d.Params, d.Ctx); len(values) > 0 {
		for key, value := range values {
			if value != "" {
				d.Conditions[key] = value
			}
		}
	}
	if values := gins.URLDefaultQuery(d.Querys, d.Ctx); len(values) > 0 {
		for key, value := range values {
			if value != "" {
				d.Conditions[key] = value
			}
		}
	}
	for k, v := range d.Conditions {
		if _, ok := jsons[k]; ok {
			if _, ok := gormColumns[k]; ok {
				d.DBWheres = append(d.DBWheres, gorms.Where(k, " IN (?) ", v))
			}
		}
	}
	return
}

func (d *DELETE) MysqlHandler() error {
	if count := len(d.API.ValidatorFuncs); count > 0 {
		for i := 0; i < count; i++ {
			if ok, err := d.API.ValidatorFuncs[i](&d.API); !ok {
				return err
			}
		}
	}
	if count, err := gorms.TotalCount(&d.Mysql, d.DBObject, nil); err != nil && err != gorm.ErrRecordNotFound {
		d.Ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"response": Resp("DB Error", err.Error())})
		return err
	} else if count <= 0 || err == gorm.ErrRecordNotFound {
		d.Ctx.JSON(http.StatusOK, gin.H{"response": Resp("No Record", []string{})})
		return nil
	}
	if err := gorms.ScopesDelete(&d.Mysql, d.DBWheres, d.DBObject); err != nil {
		d.Ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"response": Resp("DB Error", err.Error())})
		return err
	}
	d.Ctx.JSON(http.StatusOK, gin.H{"response": Resp("Success", "Record Deleted")})
	return nil
}
