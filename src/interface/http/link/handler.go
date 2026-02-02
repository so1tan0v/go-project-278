package link

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"link-service/src/domain/link"
	linkusecase "link-service/src/usecase/link"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	var rng *link.Range
	var res []linkusecase.LinkDTO
	var err error

	rngString := c.Query("range")
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
		res, err = h.useCase.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
	}

	response := make([]LinkResponse, 0, len(res))
	for _, l := range res {
		response = append(response, mapToResponse(l))
	}

	total, err := h.useCase.Count(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	start := 0
	if rng != nil {
		start = rng.Start
	}
	end := start + len(response) - 1
	if len(response) == 0 {
		end = start
	}
	c.Header("Content-Range", fmt.Sprintf("links %d-%d/%d", start, end, total))

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
		writeBindError(c, err)
		return
	}

	res, err := h.useCase.Create(c.Request.Context(), linkusecase.CreateInput{
		OriginalURL: req.OriginalURL,
		ShortName:   req.ShortName,
	})

	if err != nil {
		switch err {
		case linkusecase.ErrInvalidInput:
			c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{
				"original_url": "invalid input",
			}})
			return
		case linkusecase.ErrShortNameConflict:
			c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{
				"short_name": "short name already in use",
			}})
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
		writeBindError(c, err)

		return
	}

	res, err := h.useCase.Update(c.Request.Context(), id, linkusecase.UpdateInput{
		OriginalURL: req.OriginalURL,
		ShortName:   req.ShortName,
	})
	if err != nil {
		switch err {
		case linkusecase.ErrInvalidInput:
			c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{
				"original_url": "invalid input",
			}})

			return
		case linkusecase.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})

			return
		case linkusecase.ErrShortNameConflict:
			c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{
				"short_name": "short name already in use",
			}})

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

func writeBindError(c *gin.Context, err error) {
	var syntaxErr *json.SyntaxError
	var unmarshalTypeErr *json.UnmarshalTypeError

	if errors.Is(err, io.EOF) || errors.As(err, &syntaxErr) || errors.As(err, &unmarshalTypeErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})

		return
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make(map[string]string, len(ve))

		for _, fe := range ve {
			out[toSnakeCase(fe.Field())] = fe.Error()
		}

		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": out})

		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
}

func toSnakeCase(s string) string {
	if s == "" {
		return s
	}
	var b []rune
	runes := []rune(s)
	for i, r := range runes {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				prev := runes[i-1]
				var next rune
				if i+1 < len(runes) {
					next = runes[i+1]
				}
				prevIsLower := prev >= 'a' && prev <= 'z'
				nextIsLower := next >= 'a' && next <= 'z'
				if prevIsLower || nextIsLower {
					b = append(b, '_')
				}
			}

			b = append(b, r+('a'-'A'))

			continue
		}
		b = append(b, r)
	}

	return string(b)
}
