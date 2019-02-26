package gorms

import (
	"startkit/library/times"
	"startkit/starter"

	"github.com/jinzhu/gorm"
)

func BatchUpdates(mysql *starter.Mysql, values interface{}) (err error) {
	defer mysql.Connector()()
	return mysql.DB.Model(values).Updates(values).Error
}

func ScopesUpdate(mysql *starter.Mysql, scopes []func(*gorm.DB) *gorm.DB, obj interface{}) (err error) {
	defer mysql.Connector()()
	err = mysql.DB.Debug().Scopes(scopes...).Where(map[string]interface{}{"deleted_at": times.Zero()}).Update(obj).Error
	return
}
