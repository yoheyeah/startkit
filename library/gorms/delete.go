package gorms

import (
	"startkit/starter"

	"github.com/jinzhu/gorm"
)

func ScopesDelete(mysql *starter.Mysql, scopes []func(*gorm.DB) *gorm.DB, obj interface{}) (err error) {
	defer mysql.Connector()()
	err = mysql.DB. /* .Set("gorm:delete_option", "OPTION (OPTIMIZE FOR UNKNOWN)") */ Scopes(scopes...).Delete(obj).Error
	return
}
