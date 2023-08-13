package account

import (
	"context"
	"database/sql"

	"fiangumilar.id/e-wallet/domain"
)

type repository struct {
	db *sql.DB
}

func NewRepository(con *sql.DB) domain.AccountRepository {
	return &repository{db: con}
}

// FindByAccount implements domain.AccountRepository.
func (r repository) FindByAccount(ctx context.Context, account string) (acc domain.Account, err error) {
	query := `SELECT * FROM accounts WHERE account = ?`
	row := r.db.QueryRowContext(ctx, query, account)
	err = row.Scan(&acc.ID, &acc.UserID, &acc.Account, &acc.Balance)
	if err != nil {
		return acc, nil
	}
	return
}

// FindByUserID implements domain.AccountRepository.
func (r repository) FindByUserID(ctx context.Context, id int64) (account domain.Account, err error) {
	query := `SELECT * FROM accounts WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	err = row.Scan(&account.ID, &account.UserID, &account.Account, &account.Balance)
	if err != nil {
		return account, nil
	}
	return
}

// Update implements domain.AccountRepository.
func (r repository) Update(ctx context.Context, account *domain.Account) error {
	query := `UPDATE accounts SET account = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, &account.Account, &account.ID)
	return err
}
