package interfaces

import "github.com/gin-gonic/gin"

type AppInterface interface {
	InitRoutes(router *gin.Engine)
}
