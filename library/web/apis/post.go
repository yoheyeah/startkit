package apis

import (
	"net/http"
	"reflect"
	"startkit/library/gorms"

	"github.com/gin-gonic/gin"
)

type POST struct {
	API
	ValidatorFuncs []func() bool
	// Conditions map[string]string
	// DBWheres   []func(db *gorm.DB) *gorm.DB
}

func (p *POST) Run() (err error) {
	ptr := reflect.ValueOf(p.DBResult).Elem()
	// set the pointer
	ptr.Set(reflect.Zero(ptr.Type()))
	for _, name := range p.API.Context.App.InUseService {
		switch name {
		case "Mysql":
			err = p.MysqlHandler()
		}
	}
	return
}

func (p *POST) MysqlHandler() error {
	if err := gorms.BatchInsert(&p.Mysql, p.DBCreate); err != nil {
		p.Ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"response": Resp("DB Error", err.Error())})
		return err
	}
	p.Ctx.JSON(http.StatusOK, gin.H{"response": Resp("Success", p.DBResult)})
	return nil
}
