package types

import (
	"strconv"
)

func ToInt(value interface{}) (int, bool) {
	switch value.(type) {
	case string:
		number, _ := strconv.Atoi(value.(string))
		return number, true
	case float64:
		number := int(value.(float64))
		return number, true
	default:
		number, ok := value.(int)
		return number, ok
	}
}
