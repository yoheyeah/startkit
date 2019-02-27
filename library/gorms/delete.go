package gorms

import (
	"startkit/library/times"
	"startkit/starter"
	"time"

	"github.com/jinzhu/gorm"
)

func ScopesDelete(mysql *starter.Mysql, scopes []func(*gorm.DB) *gorm.DB, obj interface{}) (err error) {
	defer mysql.Connector()()
	err = mysql.DB.Debug().Model(obj).Scopes(scopes...).Where(map[string]interface{}{"deleted_at": times.Zero()}).Update("deleted_at", time.Now()).Error
	return
}
