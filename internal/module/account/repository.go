package account

import (
	"context"

	"fiangumilar.id/e-wallet/domain"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewAccountRepository(con *gorm.DB) domain.AccountRepository {
	return &repository{db: con}
}

// FindByUserID implements domain.AccountRepository.
func (r repository) FindByUserID(ctx context.Context, id int64) (account domain.Account, err error) {
	dataset := r.db.Debug().WithContext(ctx).Where("id = ?", id).First(&account)
	if dataset.Error != nil {
		return account, dataset.Error
	}
	return
}

// FindByAccountNumber implements domain.AccountRepository.
func (r repository) FindByAccountNumber(ctx context.Context, accNumber string) (account domain.Account, err error) {
	dataset := r.db.Debug().WithContext(ctx).Where("account_number = ?", accNumber).First(&account)
	if dataset.Error != nil {
		return account, dataset.Error
	}
	return
}

// Update implements domain.AccountRepository.
func (r repository) Update(ctx context.Context, account *domain.Account) error {
	executor := r.db.Debug().
		WithContext(ctx).
		Model(account). // Use &account if update all accounts
		Where("id = ?", account.ID).
		Update("account_number", account.AccountNumber)

	if executor.Error != nil {
		return executor.Error
	}
	return nil
}
