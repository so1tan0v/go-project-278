package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"

	"link-service/src/domain/entity"
	domain "link-service/src/domain/link"
	"link-service/src/infrastructure/database/sqlcdb"
)

/*Репозиторий для работы с PostgreSQL*/
type Repository struct {
	q *sqlcdb.Queries /*Queries для работы с базой данных*/
}

/*Метод создания нового репозитория*/
func New(db *sql.DB) *Repository {
	return &Repository{q: sqlcdb.New(db)}
}

/*Метод получения списка ссылок*/
func (r *Repository) List(ctx context.Context) ([]entity.Link, error) {
	rows, err := r.q.ListLinks(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]entity.Link, 0, len(rows))
	for _, row := range rows {
		res = append(res, fromSQLC(row))
	}

	return res, nil
}

/*Метод получения списка ссылок*/
func (r *Repository) ListWithRange(ctx context.Context, rng *domain.Range) ([]entity.Link, error) {
	rows, err := r.q.ListLinksWithRange(ctx, sqlcdb.ListLinksWithRangeParams{
		Limit:  int32(rng.End - rng.Start + 1),
		Offset: int32(rng.Start),
	})
	if err != nil {
		return nil, err
	}

	res := make([]entity.Link, 0, len(rows))
	for _, row := range rows {
		res = append(res, fromSQLC(row))
	}

	return res, nil
}

/*Метод получения общего количества ссылок*/
func (r *Repository) Count(ctx context.Context) (int64, error) {
	return r.q.CountLinks(ctx)
}

/*Метод получения ссылки по идентификатору*/
func (r *Repository) Get(ctx context.Context, id int64) (entity.Link, error) {
	row, err := r.q.GetLink(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Link{}, domain.ErrNotFound
		}
		return entity.Link{}, err
	}

	return fromSQLC(row), nil
}

/*Метод получения ссылки по short_name*/
func (r *Repository) GetByShortName(ctx context.Context, shortName string) (entity.Link, error) {
	row, err := r.q.GetLinkByShortName(ctx, shortName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Link{}, domain.ErrNotFound
		}
		return entity.Link{}, err
	}

	return fromSQLC(row), nil
}

/*Метод создания новой ссылки*/
func (r *Repository) Create(ctx context.Context, originalURL, shortName string) (entity.Link, error) {
	row, err := r.q.CreateLink(ctx, sqlcdb.CreateLinkParams{
		OriginalUrl: originalURL,
		ShortName:   shortName,
	})

	if err != nil {
		if isUniqueViolation(err) {
			return entity.Link{}, domain.ErrShortNameConflict
		}
		return entity.Link{}, err
	}

	return fromSQLC(row), nil
}

/*Метод обновления ссылки*/
func (r *Repository) Update(ctx context.Context, id int64, originalURL, shortName string) (entity.Link, error) {
	row, err := r.q.UpdateLink(ctx, sqlcdb.UpdateLinkParams{
		ID:          id,
		OriginalUrl: originalURL,
		ShortName:   shortName,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Link{}, domain.ErrNotFound
		}

		if isUniqueViolation(err) {
			return entity.Link{}, domain.ErrShortNameConflict
		}

		return entity.Link{}, err
	}

	return fromSQLC(row), nil
}

/*Метод удаления ссылки*/
func (r *Repository) Delete(ctx context.Context, id int64) error {
	_, err := r.q.DeleteLink(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}
	return nil
}

/*Метод преобразования из sqlcdb.Link в entity.Link*/
func fromSQLC(l sqlcdb.Link) entity.Link {
	return entity.Link{
		ID:          l.ID,
		OriginalURL: l.OriginalUrl,
		ShortName:   l.ShortName,
		CreatedAt:   l.CreatedAt,
	}
}

/*Метод проверки на уникальность*/
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}

/*Интерфейс для работы с базой данных*/
var _ domain.Repository = (*Repository)(nil)
