package httpinterface

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"link-service/src/domain/link"
	linkusecase "link-service/src/usecase/link"
	linkvisitusecase "link-service/src/usecase/linkvisit"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

type stubLinkUC struct {
	getByShortName func(ctx context.Context, shortName string) (linkusecase.LinkDTO, error)
}

func (s stubLinkUC) List(ctx context.Context) ([]linkusecase.LinkDTO, error)               { return nil, nil }
func (s stubLinkUC) ListWithRange(ctx context.Context, rng *link.Range) ([]linkusecase.LinkDTO, error) {
	return nil, nil
}
func (s stubLinkUC) Count(ctx context.Context) (int64, error)                          { return 0, nil }
func (s stubLinkUC) Get(ctx context.Context, id int64) (linkusecase.LinkDTO, error)    { return linkusecase.LinkDTO{}, nil }
func (s stubLinkUC) Create(ctx context.Context, in linkusecase.CreateInput) (linkusecase.LinkDTO, error) {
	return linkusecase.LinkDTO{}, nil
}
func (s stubLinkUC) Update(ctx context.Context, id int64, in linkusecase.UpdateInput) (linkusecase.LinkDTO, error) {
	return linkusecase.LinkDTO{}, nil
}
func (s stubLinkUC) Delete(ctx context.Context, id int64) error { return nil }
func (s stubLinkUC) GetByShortName(ctx context.Context, shortName string) (linkusecase.LinkDTO, error) {
	return s.getByShortName(ctx, shortName)
}

type stubVisitUC struct {
	create        func(ctx context.Context, in linkvisitusecase.CreateInput) (linkvisitusecase.LinkVisitDTO, error)
	listWithRange func(ctx context.Context, rng *link.Range) ([]linkvisitusecase.LinkVisitDTO, error)
	count         func(ctx context.Context) (int64, error)
}

func (s stubVisitUC) Create(ctx context.Context, in linkvisitusecase.CreateInput) (linkvisitusecase.LinkVisitDTO, error) {
	return s.create(ctx, in)
}
func (s stubVisitUC) ListWithRange(ctx context.Context, rng *link.Range) ([]linkvisitusecase.LinkVisitDTO, error) {
	return s.listWithRange(ctx, rng)
}
func (s stubVisitUC) Count(ctx context.Context) (int64, error) { return s.count(ctx) }

func TestRedirectCreatesVisitAndRedirects(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var created *linkvisitusecase.CreateInput

	router := gin.New()
	router.TrustedPlatform = gin.PlatformCloudflare

	linkUC := stubLinkUC{
		getByShortName: func(ctx context.Context, shortName string) (linkusecase.LinkDTO, error) {
			assert.Equal(t, "abc", shortName)
			return linkusecase.LinkDTO{
				ID:          1,
				OriginalURL: "https://example.com",
				ShortName:   "abc",
				ShortURL:    "http://localhost/r/abc",
			}, nil
		},
	}

	visitUC := stubVisitUC{
		create: func(ctx context.Context, in linkvisitusecase.CreateInput) (linkvisitusecase.LinkVisitDTO, error) {
			tmp := in
			created = &tmp
			return linkvisitusecase.LinkVisitDTO{ID: 5}, nil
		},
		listWithRange: func(ctx context.Context, rng *link.Range) ([]linkvisitusecase.LinkVisitDTO, error) {
			return nil, nil
		},
		count: func(ctx context.Context) (int64, error) { return 0, nil },
	}

	InitRoutes(router, Deps{Link: linkUC, LinkVisit: visitUC})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/r/abc", nil)
	req.Header.Set("CF-Connecting-IP", "1.2.3.4")
	req.Header.Set("User-Agent", "curl/8.5.0")
	req.Header.Set("Referer", "https://ref.example/")
	req.RemoteAddr = "10.0.0.1:1234"

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "https://example.com", w.Header().Get("Location"))

	if created == nil {
		t.Fatalf("expected visit to be created")
	}
	assert.Equal(t, int64(1), created.LinkID)
	assert.Equal(t, "1.2.3.4", created.IP)
	assert.Equal(t, "curl/8.5.0", created.UserAgent)
	assert.Equal(t, "https://ref.example/", created.Referer)
	assert.Equal(t, http.StatusFound, created.Status)
}

func TestLinkVisitsListPaginationSetsContentRange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	visits := make([]linkvisitusecase.LinkVisitDTO, 0, 11)
	now := time.Date(2025, 10, 31, 13, 1, 43, 0, time.UTC)
	for i := 0; i < 11; i++ {
		visits = append(visits, linkvisitusecase.LinkVisitDTO{
			ID:        int64(5 + i),
			LinkID:    1,
			CreatedAt: now,
			IP:        "172.18.0.1",
			UserAgent: "curl/8.5.0",
			Referer:   "",
			Status:    302,
		})
	}

	linkUC := stubLinkUC{getByShortName: func(ctx context.Context, shortName string) (linkusecase.LinkDTO, error) {
		return linkusecase.LinkDTO{}, nil
	}}

	visitUC := stubVisitUC{
		create: func(ctx context.Context, in linkvisitusecase.CreateInput) (linkvisitusecase.LinkVisitDTO, error) {
			return linkvisitusecase.LinkVisitDTO{}, nil
		},
		listWithRange: func(ctx context.Context, rng *link.Range) ([]linkvisitusecase.LinkVisitDTO, error) {
			assert.Equal(t, 10, rng.Start)
			assert.Equal(t, 20, rng.End)
			return visits, nil
		},
		count: func(ctx context.Context) (int64, error) { return 357, nil },
	}

	InitRoutes(router, Deps{Link: linkUC, LinkVisit: visitUC})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/link_visits", nil)
	req.URL.RawQuery = "range=%5B10,20%5D"

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "link_visits 10-20/357", w.Header().Get("Content-Range"))

	var got []map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, nil, err)
	assert.Equal(t, 11, len(got))
}

