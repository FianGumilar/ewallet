package transaction

import (
	"context"

	"fiangumilar.id/e-wallet/domain"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewTransactionRepository(con *gorm.DB) domain.TransactionRepository {
	return &repository{db: con}
}

// Insert implements domain.TransactionRepository.
func (r repository) Insert(ctx context.Context, transaction *domain.Transaction) error {
	executor := r.db.Debug().WithContext(ctx).Create(transaction)
	if executor.Error != nil {
		return executor.Error
	}
	return nil
}
