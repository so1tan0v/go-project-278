package linkvisit

import (
	"context"

	"link-service/src/domain/entity"
	"link-service/src/domain/link"
)

/*Репозиторий для посещений ссылок*/
type Repository interface {
	/*Создание записи посещения*/
	Create(ctx context.Context, in CreateInput) (entity.LinkVisit, error)
	/*Список посещений с range*/
	ListWithRange(ctx context.Context, rng *link.Range) ([]entity.LinkVisit, error)
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

