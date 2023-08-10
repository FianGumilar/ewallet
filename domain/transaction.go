package domain

import (
	"context"
	"time"

	"fiangumilar.id/e-wallet/dto"
)

type Transaction struct {
	ID                  int64     `gorm:"primaryKey"`
	AccountID           int64     `gorm:"column:account_id"`
	SofNumber           string    `gorm:"column:sof_number"`
	DofNumber           string    `gorm:"column:dof_number"`
	TransactionType     string    `gorm:"column:transaction_type"`
	Amount              float64   `gorm:"column:amount"`
	TransactionDateTime time.Time `gorm:"column:transaction_date_time"`
}

// SofNumber = Source of Fund (Sumber dana)
// DofNumber = Destination of Fun (Tujuan dana)

type TransactionRepository interface {
	Insert(ctx context.Context, transaction *Transaction) error
}

type TransactionService interface {
	TransaferInquiry(ctx context.Context, req dto.TransferInquiryReq) (dto.TransferInquiryRes, error)
	TransferExecute(ctx context.Context, req dto.TransferExecuteReq) error
}
