package gorms

import (
	"github.com/go-sql-driver/mysql"
)

var (
	MysqlError = map[int]string{
		1064: "Syntax Error",
		1175: "Safe Update",
		1067: "This is probably related to TIMESTAMP defaults, which have changed over time. See TIMESTAMP defaults in the Dates & Times page. (which does not exist yet)",
		1292: "DOUBLE/Integer Check for letters or other syntax errors. Check that the columns align; perhaps you think you are putting into a VARCHAR but it is aligned with a numeric column. DATETIME Check for too far in past or future. Check for between 2am and 3am on a morning when Daylight savings changed. Check for bad syntax, such as +00 timezone stuff. VARIABLE Check the allowed values for the VARIABLE you are trying to SET. LOAD DATA Look at the line that is 'bad'. Check the escape symbols, etc. Look at the datatypes.",
		1366: "DOUBLE/Integer Check for letters or other syntax errors. Check that the columns align; perhaps you think you are putting into a VARCHAR but it is aligned with a numeric column.",
		1411: "STR_TO_DATE Incorrectly formatted date",
	}
)

func IsDuplicateError(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return mysqlErr.Number == 1062
	}
	return false
}
