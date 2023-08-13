package domain

import "context"

type Account struct {
	ID      int64   `db:"id"`
	UserID  int64   `db:"user_id"`
	Account string  `db:"account"`
	Balance float64 `db:"balance"`
}

type AccountRepository interface {
	FindByUserID(ctx context.Context, id int64) (Account, error)
	FindByAccount(ctx context.Context, account string) (Account, error)
	Update(ctx context.Context, account *Account) error
}
