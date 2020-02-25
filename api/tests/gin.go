package tests

import (
	"github.com/gin-gonic/gin"
)

// SetGinContext sets custom Gin contexts
func SetGinContext(r *gin.Engine, ctxs map[string]interface{}) {
	r.Use(func(c *gin.Context) {
		for k, v := range ctxs {
			c.Set(k, v)
		}
	})
}
