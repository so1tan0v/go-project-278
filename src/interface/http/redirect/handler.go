package redirect

import (
	"net/http"

	"github.com/gin-gonic/gin"

	linkusecase "link-service/src/usecase/link"
	linkvisitusecase "link-service/src/usecase/linkvisit"
)

/*Хендлер для редиректа по короткой ссылке*/
type Handler struct {
	linkUseCase      linkusecase.UseCase
	linkVisitUseCase linkvisitusecase.UseCase
}

/*Метод создания нового хендлера*/
func NewHandler(linkUC linkusecase.UseCase, visitUC linkvisitusecase.UseCase) *Handler {
	return &Handler{
		linkUseCase:      linkUC,
		linkVisitUseCase: visitUC,
	}
}

/*Редирект по short_name с записью посещения*/
func (h *Handler) Redirect(c *gin.Context) {
	code := c.Param("code")

	l, err := h.linkUseCase.GetByShortName(c.Request.Context(), code)
	if err != nil {
		// в usecase уже нормализованы ErrNotFound, но тут неважно — отдаём 404
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	status := http.StatusFound
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	referer := c.GetHeader("Referer")

	if _, err := h.linkVisitUseCase.Create(c.Request.Context(), linkvisitusecase.CreateInput{
		LinkID:    l.ID,
		IP:        ip,
		UserAgent: userAgent,
		Referer:   referer,
		Status:    status,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.Redirect(status, l.OriginalURL)
}

