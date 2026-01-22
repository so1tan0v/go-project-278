package link

import "github.com/gin-gonic/gin"

func InitRoutes(router *gin.Engine) {
	apiRoute := router.Group("/api")

	initGetRoutes(apiRoute)
	initPostRoutes(apiRoute)
	initPutRoutes(apiRoute)
	initDeleteRoutes(apiRoute)
}

func initGetRoutes(router *gin.RouterGroup) {
	router.GET("/links", func(c *gin.Context) {

	})
	router.GET("/links/:id", func(c *gin.Context) {
		//id := c.Param("id")

	})
}

func initPostRoutes(router *gin.RouterGroup) {
	router.POST("/links", func(c *gin.Context) {

	})
}

func initDeleteRoutes(router *gin.RouterGroup) {
	router.DELETE("/links/:id", func(c *gin.Context) {
		//id := c.Param("id")
	})
}

func initPutRoutes(router *gin.RouterGroup) {
	router.PUT("/links/:id", func(c *gin.Context) {
		//id := c.Param("id")
	})
}
