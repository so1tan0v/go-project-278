package linkvisitusecase

import (
	"context"

	"link-service/src/domain/entity"
	domain "link-service/src/domain/linkvisit"
	"link-service/src/domain/link"
)

/*Сервис для работы с посещениями ссылок*/
type Service struct {
	repo domain.Repository
}

/*Метод создания нового сервиса*/
func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

/*Создание посещения*/
func (s *Service) Create(ctx context.Context, in CreateInput) (LinkVisitDTO, error) {
	v, err := s.repo.Create(ctx, domain.CreateInput{
		LinkID:    in.LinkID,
		IP:        in.IP,
		UserAgent: in.UserAgent,
		Referer:   in.Referer,
		Status:    in.Status,
	})
	if err != nil {
		return LinkVisitDTO{}, err
	}
	return toDTO(v), nil
}

/*Список посещений с range*/
func (s *Service) ListWithRange(ctx context.Context, rng *link.Range) ([]LinkVisitDTO, error) {
	visits, err := s.repo.ListWithRange(ctx, rng)
	if err != nil {
		return nil, err
	}

	res := make([]LinkVisitDTO, 0, len(visits))
	for _, v := range visits {
		res = append(res, toDTO(v))
	}
	return res, nil
}

/*Общее количество посещений*/
func (s *Service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func toDTO(v entity.LinkVisit) LinkVisitDTO {
	return LinkVisitDTO{
		ID:        v.ID,
		LinkID:    v.LinkID,
		CreatedAt: v.CreatedAt,
		IP:        v.IP,
		UserAgent: v.UserAgent,
		Referer:   v.Referer,
		Status:    v.Status,
	}
}

var _ UseCase = (*Service)(nil)

