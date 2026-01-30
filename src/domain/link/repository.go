package link

import (
	"context"

	"link-service/src/domain/entity"
)

/*Репозиторий для ссылок*/
type Repository interface {
	/*Список ссылок*/
	List(ctx context.Context) ([]entity.Link, error)
	/*Список ссылок с range*/
	ListWithRange(ctx context.Context, rng *Range) ([]entity.Link, error)
	/*Общее количество ссылок*/
	Count(ctx context.Context) (int64, error)
	/*Получение ссылки по идентификатору*/
	Get(ctx context.Context, id int64) (entity.Link, error)
	/*Создание ссылки*/
	Create(ctx context.Context, originalURL, shortName string) (entity.Link, error)
	/*Обновление ссылки*/
	Update(ctx context.Context, id int64, originalURL, shortName string) (entity.Link, error)
	/*Удаление ссылки*/
	Delete(ctx context.Context, id int64) error
}
