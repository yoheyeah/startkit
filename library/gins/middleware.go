package gins

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stvp/rollbar"
	validator "gopkg.in/go-playground/validator.v8"
)

func GinErrors() gin.HandlerFunc {
	type Resp struct {
		Message interface{} `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				switch e.Type {
				case gin.ErrorTypePublic:
					if !c.Writer.Written() {
						c.JSON(c.Writer.Status(), gin.H{
							"error_message": Resp{
								Message: "Public Error",
								Data:    e.Error(),
							},
						})
					}
				case gin.ErrorTypeBind:
					errs := e.Err.(validator.ValidationErrors)
					list := make(map[string]string)
					for _, err := range errs {
						list[err.Field] = ValidationErrorToText(err)
					}
					status := http.StatusBadRequest
					if c.Writer.Status() != http.StatusOK {
						status = c.Writer.Status()
					}
					c.JSON(status, gin.H{
						"error_message": Resp{
							Message: "Request Binding Error",
							Data:    list,
						},
					})

				default:
					rollbar.RequestError(rollbar.ERR, c.Request, e.Err)
				}

			}
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error_message": Resp{
						Message: "Internal Server Error",
						Data:    c.Errors[len(c.Errors)-1].Error(),
					},
				})
			}
		}
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Content-Disposition")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		c.Next()
	}
}