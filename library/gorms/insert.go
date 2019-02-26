package gorms

import (
	"errors"
	"fmt"
	"startkit/starter"
	"strings"
	"time"
)

func BatchInsert(mysql *starter.Mysql, values []interface{}) (err error) {
	defer mysql.Connector()()
	if len(values) == 0 {
		return errors.New("Insert values with 0 length")
	}
	var (
		mainObj    = values[0]
		mainScope  = mysql.DB.NewScope(mainObj)
		mainFields = mainScope.Fields()
		quoted     = make([]string, 0, len(mainFields))
	)
	for i := range mainFields {
		if (mainFields[i].IsPrimaryKey && mainFields[i].IsBlank) || (mainFields[i].IsIgnored) {
			continue
		}
		if mainFields[i].Field.Type().String() == "time.Time" || mainFields[i].Field.Type().String() == "*time.Time" {
			if mainFields[i].Name != "CreatedAt" && mainFields[i].Name != "UpdatedAt" {
				if mainFields[i].Field.IsValid() {
					continue
				}
			}
		}
		quoted = append(quoted, mainScope.Quote(mainFields[i].DBName))
	}
	placeholdersArr := make([]string, 0, len(values))
	for _, val := range values {
		var (
			scope        = mysql.DB.NewScope(val)
			fields       = scope.Fields()
			placeholders = make([]string, 0, len(fields))
		)
		for i := range fields {
			if (fields[i].IsPrimaryKey && fields[i].IsBlank) || (fields[i].IsIgnored) {
				continue
			}
			if fields[i].Field.Type().String() == "time.Time" || fields[i].Field.Type().String() == "*time.Time" {
				if fields[i].Name == "CreatedAt" || fields[i].Name == "UpdatedAt" {
					placeholders = append(placeholders, scope.AddToVars(time.Now().Format("2006-01-02 15:04:05")))
				}
				if fields[i].Field.IsValid() {
					continue
				}
			} else {
				placeholders = append(placeholders, scope.AddToVars(fields[i].Field.Interface()))
			}
		}
		placeholdersStr := "(" + strings.Join(placeholders, ", ") + ")"
		placeholdersArr = append(placeholdersArr, placeholdersStr)
		mainScope.SQLVars = append(mainScope.SQLVars, scope.SQLVars...)
	}
	mainScope.Raw(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		mainScope.QuotedTableName(),
		strings.Join(quoted, ", "),
		strings.Join(placeholdersArr, ", "),
	))
	if _, err := mainScope.SQLDB().Exec(mainScope.SQL, mainScope.SQLVars...); err != nil {
		mysql.DB.Error = err
		return err
	}
	return err
}
