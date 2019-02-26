package apis

import (
	"net/http"
	"startkit/library/gorms"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type PUT struct {
	API
	Conditions map[string]string
	DBWheres   []func(db *gorm.DB) *gorm.DB
}

func (p *PUT) Run() (err error) {
	p.FillInIDer()
	for _, name := range p.API.Context.App.InUseService {
		switch name {
		case "Mysql":
			err = p.MysqlHandler()
		}
	}
	return
}

func (p *PUT) FillInIDer() {
	if p.IDer != "" {
		if id := p.Ctx.Param(p.IDer); id != "" {
			p.DBWheres = append(p.DBWheres, gorms.Where(p.IDer, " = ? ", id))
		}
	}
	return
}

func (p *PUT) MysqlHandler() error {
	if count, err := gorms.TotalCount(&p.Mysql, p.DBObject, nil); err != nil && err != gorm.ErrRecordNotFound {
		p.Ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"response": Resp("DB Error", err.Error())})
		return err
	} else if count <= 0 || err == gorm.ErrRecordNotFound {
		p.Ctx.JSON(http.StatusOK, gin.H{"response": Resp("No Record", map[string]int{"count": count})})
		return nil
	}
	if err := gorms.ScopesUpdate(&p.Mysql, p.DBWheres, p.DBUpdate); err != nil {
		p.Ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"response": Resp("DB Error", err.Error())})
		return err
	}
	p.Ctx.JSON(http.StatusOK, gin.H{"response": Resp("Success", p.DBUpdate)})
	return nil
}
