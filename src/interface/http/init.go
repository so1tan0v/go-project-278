package httpinterface

import (
	"net/http"

	"link-service/src/interface/http/link"
	"link-service/src/interface/http/ping"
	linkusecase "link-service/src/usecase/link"

	"github.com/gin-gonic/gin"
)

/*Зависимости для инициализации маршрутов*/
type Deps struct {
	Link linkusecase.UseCase
}

/*Метод инициализации маршрутов*/
func InitRoutes(router *gin.Engine, deps Deps) {
	ping.RegisterRoutes(router)

	apiRoute := router.Group("/api")
	linkHandler := link.NewHandler(deps.Link)
	link.RegisterRoutes(apiRoute, linkHandler)

	/*Метод обработки не найденных маршрутов*/
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})
}
