package gorms

import (
	"startkit/starter"
)

func Create(mysql *starter.Mysql, obj interface{}, where interface{}) (err error) {
	defer mysql.Connector()()
	return mysql.DB.FirstOrCreate(obj, where).Error
}
