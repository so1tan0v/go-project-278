package link

import (
	"github.com/gin-gonic/gin"
)

/*Метод регистрации маршрутов*/
func RegisterRoutes(router *gin.RouterGroup, h *Handler) {
	router.GET("/links", h.List)          /*Маршрут для получения списка ссылок*/
	router.GET("/links/:id", h.Get)       /*Маршрут для получения ссылки по идентификатору*/
	router.POST("/links", h.Create)       /*Маршрут для создания ссылки*/
	router.PUT("/links/:id", h.Update)    /*Маршрут для обновления ссылки*/
	router.DELETE("/links/:id", h.Delete) /*Маршрут для удаления ссылки*/
}
