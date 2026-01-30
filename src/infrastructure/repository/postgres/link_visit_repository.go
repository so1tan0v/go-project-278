package postgres

import (
	"context"
	"database/sql"

	"link-service/src/domain/entity"
	"link-service/src/domain/link"
	domain "link-service/src/domain/linkvisit"
	"link-service/src/infrastructure/database/sqlcdb"
)

/*Репозиторий посещений ссылок для PostgreSQL*/
type LinkVisitRepository struct {
	q *sqlcdb.Queries
}

/*Метод создания нового репозитория посещений*/
func NewLinkVisitRepository(db *sql.DB) *LinkVisitRepository {
	return &LinkVisitRepository{q: sqlcdb.New(db)}
}

/*Создание записи посещения*/
func (r *LinkVisitRepository) Create(ctx context.Context, in domain.CreateInput) (entity.LinkVisit, error) {
	row, err := r.q.CreateLinkVisit(ctx, sqlcdb.CreateLinkVisitParams{
		LinkID:    in.LinkID,
		Ip:        in.IP,
		UserAgent: in.UserAgent,
		Referer:   in.Referer,
		Status:    int32(in.Status),
	})
	if err != nil {
		return entity.LinkVisit{}, err
	}
	return fromSQLCVisit(row), nil
}

/*Список посещений с range*/
func (r *LinkVisitRepository) ListWithRange(ctx context.Context, rng *link.Range) ([]entity.LinkVisit, error) {
	rows, err := r.q.ListLinkVisitsWithRange(ctx, sqlcdb.ListLinkVisitsWithRangeParams{
		Limit:  int32(rng.End - rng.Start + 1),
		Offset: int32(rng.Start),
	})
	if err != nil {
		return nil, err
	}

	res := make([]entity.LinkVisit, 0, len(rows))
	for _, row := range rows {
		res = append(res, fromSQLCVisit(row))
	}
	return res, nil
}

/*Общее количество посещений*/
func (r *LinkVisitRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountLinkVisits(ctx)
}

func fromSQLCVisit(v sqlcdb.LinkVisit) entity.LinkVisit {
	return entity.LinkVisit{
		ID:        v.ID,
		LinkID:    v.LinkID,
		IP:        v.Ip,
		UserAgent: v.UserAgent,
		Referer:   v.Referer,
		Status:    int(v.Status),
		CreatedAt: v.CreatedAt,
	}
}

var _ domain.Repository = (*LinkVisitRepository)(nil)

