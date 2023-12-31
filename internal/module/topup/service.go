package topup

import (
	"context"
	"fmt"
	"time"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
	"github.com/google/uuid"
)

type topUpService struct {
	notificationService   domain.NotificationService
	midtransService       domain.MidtransService
	topUpRepository       domain.TopUpRepository
	accountRepository     domain.AccountRepository
	transactionRepository domain.TransactionRepository
}

func NewTopUpService(
	notificationService domain.NotificationService,
	midtransService domain.MidtransService,
	topUpRepository domain.TopUpRepository,
	accountRepository domain.AccountRepository,
	transactionRepository domain.TransactionRepository,
) domain.TopUpService {
	return &topUpService{
		notificationService:   notificationService,
		midtransService:       midtransService,
		topUpRepository:       topUpRepository,
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
	}
}

// InitializeTopUp implements domain.TopUpService.
func (t topUpService) InitializeTopUp(ctx context.Context, req dto.TopUpRequest) (dto.TopUpResponse, error) {
	topUp := domain.TopUp{
		ID:     uuid.NewString(),
		UserID: req.UserID,
		Status: 0,
		Amount: req.Amount,
	}
	err := t.midtransService.GenerateSnapURL(ctx, &topUp)
	if err != nil {
		return dto.TopUpResponse{}, err
	}

	err = t.topUpRepository.Insert(ctx, &topUp)
	if err != nil {
		return dto.TopUpResponse{}, err
	}

	return dto.TopUpResponse{
		SnapURL: topUp.SnapURL,
	}, nil
}

// ConfirmedTopUp implements domain.TopUpService.
func (t topUpService) ConfirmedTopUp(ctx context.Context, id string) error {
	topup, err := t.topUpRepository.FindById(ctx, id)
	if err != nil {
		return domain.TopUpReqNotFound
	}

	if topup == (domain.TopUp{}) {
		return domain.TopUpReqNotFound
	}

	account, err := t.accountRepository.FindByUserID(ctx, topup.UserID)
	if err != nil {
		return domain.ErrAccountNotFound
	}

	if account == (domain.Account{}) {
		return domain.ErrAccountNotFound
	}

	err = t.transactionRepository.Insert(ctx, &domain.Transaction{
		AccountID:           account.ID,
		SofNumber:           "00",
		DofNumber:           account.Account,
		TransactionType:     "C",
		Amount:              topup.Amount,
		TransactionDateTime: time.Now(),
	})

	if err != nil {
		return err
	}

	account.Balance += topup.Amount
	err = t.accountRepository.Update(ctx, &account)

	data := map[string]string{
		"amount": fmt.Sprintf("%.2f", topup.Amount),
	}

	_ = t.notificationService.Insert(ctx, account.UserID, "TOPUP_SUCCESS", data)

	return nil
}
