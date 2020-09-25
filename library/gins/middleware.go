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
						c.AbortWithStatusJSON(c.Writer.Status(), gin.H{
							"error_message": Resp{
								Message: "Public Error",
								Data:    e.Error(),
							},
						})
					}
				case gin.ErrorTypeBind:
					var (
						list   = make(map[string]string)
						status = http.StatusBadRequest
					)
					if errs, ok := e.Err.(validator.ValidationErrors); ok {
						for _, err := range errs {
							list[err.Field] = ValidationErrorToText(err)
						}
					} else {
						list["Error"] = e.Err.Error()
					}
					if c.Writer.Status() != http.StatusOK {
						status = c.Writer.Status()
					}
					c.AbortWithStatusJSON(status, gin.H{
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
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
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
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, HEAD, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Accept, X-Requested-With, Access-Control-Request-Method, Cache-Control, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Content-Disposition, Authorization, Access-Control-Request-Headers")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		c.Next()
	}
}

func DownloadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// alone dns prefetch
		c.Writer.Header().Set("X-DNS-Prefetch-Control", "on")
		// IE No Open
		c.Writer.Header().Set("X-Download-Options", "noopen")
		// not cache
		c.Writer.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		// set Expire
		c.Writer.Header().Set("Expires", "max-age=0")
		// Content Security Policy
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		// xss protect
		// it will caught some problems is old IE
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		// Referrer Policy
		c.Writer.Header().Set("Referrer-Header", "no-referrer")
		// cros frame, allow same origin
		c.Writer.Header().Set("X-Frame-Options", "SAMEORIGIN")
		// HSTS
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=5184000;includeSubDomains")
		// no sniff
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		// set stream
		c.Writer.Header().Set("Content-type", "application/octet-stream")
		c.Next()
	}
}
