package structs

import (
	"reflect"
	"strings"
)

func GetTags(obj interface{}, tag string) (conditions map[string]string) {
	val := reflect.ValueOf(obj)
	field := ""
	conditions = make(map[string]string)
	for i := 0; i < val.Type().NumField(); i++ {
		field = val.Type().Field(i).Tag.Get(tag)
		if field == "" || field == "-" || field == "_" {
			continue
		}
		conditions[field] = ""
	}
	return conditions
}

func GetTagsValueWithSpliter(obj interface{}, tag string, spliter string) (conditions map[string]string) {
	val := reflect.ValueOf(obj)
	field := ""
	conditions = make(map[string]string)
	for i := 0; i < val.Type().NumField(); i++ {
		field = val.Type().Field(i).Tag.Get(tag)
		if field == "" || field == "-" || field == "_" {
			continue
		}
		field = strings.Split(field, spliter)[1]
		conditions[field] = ""
	}
	return conditions
}
