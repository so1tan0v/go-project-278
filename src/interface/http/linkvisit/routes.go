package linkvisit

import "github.com/gin-gonic/gin"

/*Метод регистрации маршрутов*/
func RegisterRoutes(router *gin.RouterGroup, h *Handler) {
	router.GET("/link_visits", h.List)
}

