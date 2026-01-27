package link

import (
	"net/http"
	"strconv"

	linkusecase "link-service/src/usecase/link"

	"github.com/gin-gonic/gin"
)

/*Хендлер для работы с ссылками*/
type Handler struct {
	useCase linkusecase.UseCase /*UseCase для работы с ссылками*/
}

/*Метод создания нового хендлера*/
func NewHandler(useCase linkusecase.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

/*Метод получения списка ссылок*/
func (h *Handler) List(c *gin.Context) {
	res, err := h.useCase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	response := make([]LinkResponse, 0, len(res))
	for _, l := range res {
		response = append(response, mapToResponse(l))
	}

	c.JSON(http.StatusOK, response)
}

/*Метод получения ссылки по идентификатору*/
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	res, err := h.useCase.Get(c.Request.Context(), id)
	if err != nil {
		if err == linkusecase.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, mapToResponse(res))
}

/*Метод создания новой ссылки*/
func (h *Handler) Create(c *gin.Context) {
	var req CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})

		return
	}

	res, err := h.useCase.Create(c.Request.Context(), linkusecase.CreateInput{
		OriginalURL: req.OriginalURL,
		ShortName:   req.ShortName,
	})

	if err != nil {
		switch err {
		case linkusecase.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		case linkusecase.ErrShortNameConflict:
			c.JSON(http.StatusConflict, gin.H{"error": "short_name already exists"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
	}

	c.JSON(http.StatusCreated, mapToResponse(res))
}

/*Метод обновления ссылки*/
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})

		return
	}

	var req UpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})

		return
	}

	res, err := h.useCase.Update(c.Request.Context(), id, linkusecase.UpdateInput{
		OriginalURL: req.OriginalURL,
		ShortName:   req.ShortName,
	})
	if err != nil {
		switch err {
		case linkusecase.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		case linkusecase.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		case linkusecase.ErrShortNameConflict:
			c.JSON(http.StatusConflict, gin.H{"error": "short_name already exists"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
	}

	c.JSON(http.StatusOK, mapToResponse(res))
}

/*Метод удаления ссылки*/
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})

		return
	}

	if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
		if err == linkusecase.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})

		return
	}

	c.Status(http.StatusNoContent)
}

/*Метод преобразования из linkusecase.LinkDTO в LinkResponse*/
func mapToResponse(l linkusecase.LinkDTO) LinkResponse {
	return LinkResponse{
		ID:          l.ID,
		OriginalURL: l.OriginalURL,
		ShortName:   l.ShortName,
		ShortURL:    l.ShortURL,
	}
}
