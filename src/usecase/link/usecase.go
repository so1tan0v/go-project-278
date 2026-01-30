package linkusecase

import (
	"context"
	"link-service/src/domain/link"
)

/*Интерфейс для работы с ссылками*/
type UseCase interface {
	/*Метод получения списка ссылок*/
	List(ctx context.Context) ([]LinkDTO, error)
	/*Список ссылок с range*/
	ListWithRange(ctx context.Context, rng *link.Range) ([]LinkDTO, error)
	/*Общее количество ссылок*/
	Count(ctx context.Context) (int64, error)
	/*Метод получения ссылки по идентификатору*/
	Get(ctx context.Context, id int64) (LinkDTO, error)
	/*Метод создания новой ссылки*/
	Create(ctx context.Context, in CreateInput) (LinkDTO, error)
	/*Метод обновления ссылки*/
	Update(ctx context.Context, id int64, in UpdateInput) (LinkDTO, error)
	/*Метод удаления ссылки*/
	Delete(ctx context.Context, id int64) error
}

/*DTO для создания ссылки*/
type CreateInput struct {
	OriginalURL string /*Исходная ссылка*/
	ShortName   string /*Короткая ссылка*/
}

/*DTO для обновления ссылки*/
type UpdateInput struct {
	OriginalURL string /*Исходная ссылка*/
	ShortName   string /*Короткая ссылка*/
}
