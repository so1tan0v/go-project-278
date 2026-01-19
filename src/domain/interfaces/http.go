package interfaces

import "github.com/gin-gonic/gin"

type AppInterface interface {
	Init(router *gin.Engine)
}
