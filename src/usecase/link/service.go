package linkusecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/url"
	"strings"

	"link-service/src/domain/entity"
	domain "link-service/src/domain/link"
)

/*Сервис для работы с ссылками*/
type Service struct {
	repo    domain.Repository
	baseURL string
}

/*Метод создания нового сервиса*/
func NewService(repo domain.Repository, baseURL string) *Service {
	return &Service{
		repo:    repo,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

/*Метод получения списка ссылок*/
func (s *Service) List(ctx context.Context) ([]LinkDTO, error) {
	links, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]LinkDTO, 0, len(links))
	for _, l := range links {
		res = append(res, s.toDTO(l))
	}

	return res, nil
}

/*Метод получения списка ссылок*/
func (s *Service) ListWithRange(ctx context.Context, rng *domain.Range) ([]LinkDTO, error) {
	links, err := s.repo.ListWithRange(ctx, rng)
	if err != nil {
		return nil, err
	}

	res := make([]LinkDTO, 0, len(links))
	for _, l := range links {
		res = append(res, s.toDTO(l))
	}

	return res, nil
}

/*Метод получения общего количества ссылок*/
func (s *Service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

/*Метод получения ссылки по идентификатору*/
func (s *Service) Get(ctx context.Context, id int64) (LinkDTO, error) {
	l, err := s.repo.Get(ctx, id)
	if err != nil {
		return LinkDTO{}, mapDomainError(err)
	}
	return s.toDTO(l), nil
}

/*Метод создания новой ссылки*/
func (s *Service) Create(ctx context.Context, in CreateInput) (LinkDTO, error) {
	if err := validateOriginalURL(in.OriginalURL); err != nil {
		return LinkDTO{}, ErrInvalidInput
	}

	shortName := strings.TrimSpace(in.ShortName)
	if shortName != "" {
		l, err := s.repo.Create(ctx, in.OriginalURL, shortName)
		if err != nil {
			return LinkDTO{}, mapDomainError(err)
		}

		return s.toDTO(l), nil
	}

	/*Если short_name не задан — генерируем уникальное имя.*/
	const attempts = 8
	for i := 0; i < attempts; i++ {
		gen := generateShortName(6)
		l, err := s.repo.Create(ctx, in.OriginalURL, gen)

		if err == nil {
			return s.toDTO(l), nil
		}

		if errors.Is(err, domain.ErrShortNameConflict) {
			continue
		}

		return LinkDTO{}, mapDomainError(err)
	}

	return LinkDTO{}, ErrShortNameConflict
}

/*Метод обновления ссылки*/
func (s *Service) Update(ctx context.Context, id int64, in UpdateInput) (LinkDTO, error) {
	if err := validateOriginalURL(in.OriginalURL); err != nil {
		return LinkDTO{}, ErrInvalidInput
	}

	shortName := strings.TrimSpace(in.ShortName)
	if shortName == "" {
		return LinkDTO{}, ErrInvalidInput
	}

	l, err := s.repo.Update(ctx, id, in.OriginalURL, shortName)

	if err != nil {
		return LinkDTO{}, mapDomainError(err)
	}

	return s.toDTO(l), nil
}

/*Метод удаления ссылки*/
func (s *Service) Delete(ctx context.Context, id int64) error {
	return mapDomainError(s.repo.Delete(ctx, id))
}

/*Метод преобразования из entity.Link в LinkDTO*/
func (s *Service) toDTO(l entity.Link) LinkDTO {
	return LinkDTO{
		ID:          l.ID,
		OriginalURL: l.OriginalURL,
		ShortName:   l.ShortName,
		ShortURL:    s.baseURL + "/r/" + l.ShortName,
	}
}

/*Метод валидации исходной ссылки*/
func validateOriginalURL(raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ErrInvalidInput
	}

	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ErrInvalidInput
	}

	return nil
}

/*Метод генерации короткой ссылки*/
func generateShortName(n int) string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	s := base64.RawURLEncoding.EncodeToString(b)

	if len(s) < n {
		return s
	}

	return s[:n]
}

/*Метод преобразования ошибки из domain в usecase*/
func mapDomainError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return ErrNotFound
	case errors.Is(err, domain.ErrShortNameConflict):
		return ErrShortNameConflict
	case errors.Is(err, domain.ErrInvalidInput):
		return ErrInvalidInput
	default:
		return err
	}
}

var _ UseCase = (*Service)(nil)
