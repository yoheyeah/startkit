package gins

import (
	"github.com/gin-gonic/gin"
)

func URLQuery(querys []string, c *gin.Context) (values map[string]string) {
	if count := len(querys); count > 0 {
		values = make(map[string]string)
		for i := 0; i < count; i++ {
			values[querys[i]] = c.Query(querys[i])
		}
	}
	return
}

func URLDefaultQuery(querys map[string]string, c *gin.Context) (values map[string]string) {
	if count := len(querys); count > 0 {
		values = make(map[string]string)
		for key, value := range querys {
			values[key] = c.DefaultQuery(key, value)
		}
	}
	return
}

func URLQueryMap(querys map[string]string, c *gin.Context) (values map[string]map[string]string) {
	if count := len(querys); count > 0 {
		values = make(map[string]map[string]string)
		for key := range querys {
			values[key] = make(map[string]string)
			if val, ok := c.GetQueryMap(key); ok {
				values[key] = val
			}
		}
	}
	return
}

func URLParam(params []string, c *gin.Context) (values map[string]string) {
	if count := len(params); count > 0 {
		values = make(map[string]string)
		for i := 0; i < count; i++ {
			values[params[i]] = c.Param(params[i])
		}
	}
	return
}
