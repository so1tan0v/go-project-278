package httpInterface

import (
	"link-service/src/interfaces/http/link"
	"link-service/src/interfaces/http/ping"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	ping.InitRoutes(router)

	link.InitRoutes(router)
}
