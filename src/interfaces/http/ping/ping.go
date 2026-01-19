package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	initGetHandlers(router)
}

func initGetHandlers(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}
