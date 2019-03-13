package gorms

import (
	"startkit/starter"

	"github.com/jinzhu/gorm"
)

// TotalCount get the count of the rows of table, the obj is model address value,
// maps can be map[string]interface{}{} or string for where condition of query
func TotalCount(mysql *starter.Mysql, obj, maps interface{}) (count int, err error) {
	defer mysql.Connector()()
	if maps != nil {
		err = mysql.DB.Debug().Model(obj).Where(maps). /*.Where(map[string]interface{}{"deleted_at": times.Zero()})*/ Count(&count).Error
		return
	}
	err = mysql.DB.Debug().Model(obj). /*.Where(map[string]interface{}{"deleted_at": times.Zero()})*/ Count(&count).Error
	return
}

func ScopesQuery(mysql *starter.Mysql, scopes []func(*gorm.DB) *gorm.DB, obj interface{}) (err error) {
	defer mysql.Connector()()
	err = mysql.DB.Debug().Scopes(scopes...). /*.Where(map[string]interface{}{"deleted_at": times.Zero()})*/ Find(obj).Error
	return
}
