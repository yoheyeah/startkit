package apis

import (
	"net/http"
	"startkit/library/gorms"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type DELETE struct {
	API
	Conditions map[string]string
	DBWheres   []func(db *gorm.DB) *gorm.DB
}

func (d *DELETE) Run() (err error) {
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
			d.DBWheres = append(d.DBWheres, gorms.Where(d.IDer, " = ? ", id))
		}
	}
	return
}

func (d *DELETE) MysqlHandler() error {
	if count, err := gorms.TotalCount(&d.Mysql, d.DBObject, nil); err != nil && err != gorm.ErrRecordNotFound {
		d.Ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"response": Resp("DB Error", err.Error())})
		return err
	} else if count <= 0 || err == gorm.ErrRecordNotFound {
		d.Ctx.JSON(http.StatusOK, gin.H{"response": Resp("No Record", map[string]int{"count": count})})
		return nil
	}
	if err := gorms.ScopesDelete(&d.Mysql, d.DBWheres, d.DBObject); err != nil {
		d.Ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"response": Resp("DB Error", err.Error())})
		return err
	}
	d.Ctx.JSON(http.StatusOK, gin.H{"response": Resp("Success", "Record Deleted")})
	return nil
}
