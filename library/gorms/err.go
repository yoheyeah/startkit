package gorms

import (
	"github.com/go-sql-driver/mysql"
)

func IsDuplicateError(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return mysqlErr.Number == 1062
	}
	return false
}
