package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}
