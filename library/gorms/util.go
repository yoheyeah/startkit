package gorms

import (
	"github.com/jinzhu/gorm"
)

func Error(db *gorm.DB) error {
	errs := db.GetErrors()
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
