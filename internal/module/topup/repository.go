package topup

import (
	"context"
	"database/sql"

	"fiangumilar.id/e-wallet/domain"
)

type repository struct {
	db *sql.DB
}

func NewTopUpRepository(con *sql.DB) domain.TopUpRepository {
	return &repository{db: con}
}

// FindById implements domain.TopUpRepository.
func (r repository) FindById(ctx context.Context, id string) (topup domain.TopUp, err error) {
	query := `SELECT * FROM topup WHERE id = ?`

	rows := r.db.QueryRowContext(ctx, query, id)
	err = rows.Scan(&topup.ID, &topup.UserID, &topup.Amount, &topup.Status, &topup.SnapURL)
	if err != nil {
		return topup, nil
	}
	return
}

// Insert implements domain.TopUpRepository.
func (r repository) Insert(ctx context.Context, t *domain.TopUp) error {
	query := `INSERT INTO topup (user_id, status, amount, snap_url) VALUES (?, ?, ?, ?);`

	_, err := r.db.ExecContext(ctx, query, t.UserID, t.Status, t.Amount, t.SnapURL)
	return err
}

// Update implements domain.TopUpRepository.
func (r repository) Update(ctx context.Context, t *domain.TopUp) error {
	query := `UPDATE topup WHERE id = ? SET user_id = ?, status = ?, amount = ?, snap_url = ?;`

	_, err := r.db.ExecContext(ctx, query, t.ID, t.UserID, t.Status, t.Amount, t.SnapURL)
	return err
}
