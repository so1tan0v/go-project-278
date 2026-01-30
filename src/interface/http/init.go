package httpinterface

import (
	"net/http"

	"link-service/src/interface/http/link"
	"link-service/src/interface/http/linkvisit"
	"link-service/src/interface/http/ping"
	"link-service/src/interface/http/redirect"
	linkusecase "link-service/src/usecase/link"
	linkvisitusecase "link-service/src/usecase/linkvisit"

	"github.com/gin-gonic/gin"
)

/*Зависимости для инициализации маршрутов*/
type Deps struct {
	Link      linkusecase.UseCase
	LinkVisit linkvisitusecase.UseCase
}

/*Метод инициализации маршрутов*/
func InitRoutes(router *gin.Engine, deps Deps) {
	ping.RegisterRoutes(router)

	redirectHandler := redirect.NewHandler(deps.Link, deps.LinkVisit)
	router.GET("/r/:code", redirectHandler.Redirect)

	apiRoute := router.Group("/api")
	linkHandler := link.NewHandler(deps.Link)
	link.RegisterRoutes(apiRoute, linkHandler)

	linkVisitHandler := linkvisit.NewHandler(deps.LinkVisit)
	linkvisit.RegisterRoutes(apiRoute, linkVisitHandler)

	/*Метод обработки не найденных маршрутов*/
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})
}
