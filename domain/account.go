package domain

import "context"

type Account struct {
	ID            int64   `gorm:"primaryKey"`
	UserID        int64   `gorm:"column:user_id"`
	AccountNumber string  `gorm:"column:account_number"`
	Balance       float64 `gorm:"column:balance"`
}

type AccountRepository interface {
	FindByUserID(ctx context.Context, id int64) (Account, error)
	FindByAccountNumber(ctx context.Context, accNumber string) (Account, error)
	Update(ctx context.Context, account *Account) error
}
