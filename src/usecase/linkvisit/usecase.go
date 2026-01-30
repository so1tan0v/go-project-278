package linkvisitusecase

import (
	"context"

	"link-service/src/domain/link"
)

/*Интерфейс для работы с посещениями ссылок*/
type UseCase interface {
	/*Создание посещения*/
	Create(ctx context.Context, in CreateInput) (LinkVisitDTO, error)
	/*Список посещений с range*/
	ListWithRange(ctx context.Context, rng *link.Range) ([]LinkVisitDTO, error)
	/*Общее количество посещений*/
	Count(ctx context.Context) (int64, error)
}

/*Входные параметры для создания посещения*/
type CreateInput struct {
	LinkID    int64
	IP        string
	UserAgent string
	Referer   string
	Status    int
}

