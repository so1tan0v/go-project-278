package linkvisit

import (
	"fmt"
	"net/http"

	"link-service/src/domain/link"
	linkvisitusecase "link-service/src/usecase/linkvisit"

	"github.com/gin-gonic/gin"
)

/*Хендлер для работы с посещениями ссылок*/
type Handler struct {
	useCase linkvisitusecase.UseCase
}

/*Метод создания нового хендлера*/
func NewHandler(useCase linkvisitusecase.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

/*Метод получения списка посещений*/
func (h *Handler) List(c *gin.Context) {
	rngString := c.Query("range")
	var (
		rng *link.Range
		res []linkvisitusecase.LinkVisitDTO
		err error
	)

	if rngString != "" {
		rng, err = link.ParseRange(rngString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err = h.useCase.ListWithRange(c.Request.Context(), rng)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
	} else {
		// по умолчанию — первая страница
		rng = &link.Range{Start: 0, End: 49}
		res, err = h.useCase.ListWithRange(c.Request.Context(), rng)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
	}

	total, err := h.useCase.Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	start := rng.Start
	end := start + len(res) - 1
	if len(res) == 0 {
		end = start
	}
	c.Header("Content-Range", fmt.Sprintf("link_visits %d-%d/%d", start, end, total))

	c.JSON(http.StatusOK, res)
}
