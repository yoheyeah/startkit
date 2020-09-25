package gorms

import (
	"startkit/starter"
)

func Create(mysql *starter.Mysql, obj interface{}, where interface{}) (err error) {
	defer mysql.Connector()()
	return mysql.DB.FirstOrCreate(obj, where).Error
}

func TransactionCreate(mysql *starter.Mysql, obj interface{}, where interface{}) (err error) {
	close := mysql.Connector()
	tx := mysql.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	if err := tx.FirstOrCreate(obj, where).Error; err != nil {
		tx.Rollback()
		return err
	}
	close()
	return tx.Commit().Error
}
