package factor

import (
	"context"
	"database/sql"

	"fiangumilar.id/e-wallet/domain"
)

type repository struct {
	db *sql.DB
}

func NewFactorRepository(con *sql.DB) domain.FactorRepository {
	return &repository{
		db: con,
	}
}

// FindByUser implements domain.FactorRepository.
func (r repository) FindByUser(ctx context.Context, id int64) (factor domain.Factor, err error) {
	query := `SELECT * FROM factors WHERE user_id = ?`

	rows := r.db.QueryRowContext(ctx, query, id)
	err = rows.Scan(&factor.ID, &factor.UserId, &factor.PIN)
	if err != nil {
		return factor, err
	}
	return
}
