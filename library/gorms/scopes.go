package gorms

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
)

func Pagination(start, limit string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		s, err := strconv.Atoi(start)
		if err != nil {
			return db
		}
		l, err := strconv.Atoi(limit)
		if err != nil {
			return db
		}
		if l == 0 || l == -1 {
			return db
		}
		return db.Limit(l).Offset(s)
	}
}

func Where(column, cond string, arg interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes().Where(fmt.Sprintf("%s", column+cond), arg)
	}
}

func OrderBy(field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes().Order(field)
	}
}

func IsNotNull(field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes().Where(field + " IS NOT NULL AND " + field + " NOT IN ('')")
	}
}

func IsNull(field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes().Where(field + " IS NULL OR " + field + " IN ('')")
	}
}

func Preload(column string, conditions []interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(conditions) > 0 {
			return db.Scopes().Preload(column, conditions)
		}
		return db.Scopes().Preload(column)
	}
}
