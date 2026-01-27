package linkusecase

import (
	"context"
)

// UseCase — интерфейс слоя usecase.
// HTTP слой должен зависеть только от этого интерфейса.
type UseCase interface {
	List(ctx context.Context) ([]LinkDTO, error)
	Get(ctx context.Context, id int64) (LinkDTO, error)
	Create(ctx context.Context, in CreateInput) (LinkDTO, error)
	Update(ctx context.Context, id int64, in UpdateInput) (LinkDTO, error)
	Delete(ctx context.Context, id int64) error
}

type CreateInput struct {
	OriginalURL string
	ShortName   string
}

type UpdateInput struct {
	OriginalURL string
	ShortName   string
}
