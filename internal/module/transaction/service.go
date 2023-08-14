package transaction

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
	"fiangumilar.id/e-wallet/internal/utils"
)

type service struct {
	accountRepository      domain.AccountRepository
	transactionRepository  domain.TransactionRepository
	cacheRepository        domain.CacheRepository
	notificationRepository domain.NotificationRepository
	hub                    *dto.Hub
}

func NewTransactionService(
	accountRepository domain.AccountRepository,
	transactionRepository domain.TransactionRepository,
	cacheRepository domain.CacheRepository,
	notificationRepository domain.NotificationRepository,
	hub *dto.Hub,
) domain.TransactionService {
	return &service{
		accountRepository:      accountRepository,
		transactionRepository:  transactionRepository,
		cacheRepository:        cacheRepository,
		notificationRepository: notificationRepository,
		hub:                    hub,
	}
}

// TransaferInquiry implements domain.TransactionService.
// TransferInquiry implements domain.TransactionService.
func (s service) TransferInquiry(ctx context.Context, req dto.TransferInquiryReq) (dto.TransferInquiryRes, error) {
	user := ctx.Value("x-user").(dto.UserData)

	myAccount, err := s.accountRepository.FindByUserID(ctx, user.ID)
	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if myAccount == (domain.Account{}) {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound
	}

	dofAccount, err := s.accountRepository.FindByAccount(ctx, req.Account)
	if err != nil {
		return dto.TransferInquiryRes{}, domain.ErrInquiryNotFound
	}

	if dofAccount == (domain.Account{}) {
		return dto.TransferInquiryRes{}, domain.ErrInquiryNotFound
	}
	if req.Amount > myAccount.Balance {
		return dto.TransferInquiryRes{}, domain.ErrInsufficientBalance
	}

	inquiryKey := utils.GenerateRandomString(32)

	jsonData, _ := json.Marshal(req)
	_ = s.cacheRepository.Set(inquiryKey, jsonData)

	return dto.TransferInquiryRes{
		InquiryKey: inquiryKey,
	}, nil
}

// TransferInquiryExecute implements domain.TransactionService.
func (s service) TransferExecute(ctx context.Context, req dto.TransferExecuteReq) error {
	val, err := s.cacheRepository.Get(req.InquiryKey)
	log.Printf("value: %s", val)
	if err != nil {
		return domain.ErrInquiryNotFound
	}

	var reqInq dto.TransferInquiryReq
	_ = json.Unmarshal(val, &reqInq)
	if reqInq == (dto.TransferInquiryReq{}) {
		return domain.ErrInquiryNotFound
	}

	user := ctx.Value("x-user").(dto.UserData)
	myAccount, err := s.accountRepository.FindByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	dofAccount, err := s.accountRepository.FindByAccount(ctx, reqInq.Account)
	if err != nil {
		return err
	}

	debitTransaction := domain.Transaction{
		AccountID:           myAccount.ID,
		SofNumber:           myAccount.Account,
		DofNumber:           dofAccount.Account,
		TransactionType:     "D",
		Amount:              reqInq.Amount,
		TransactionDateTime: time.Now(),
	}

	err = s.transactionRepository.Insert(ctx, &debitTransaction)
	if err != nil {
		return err
	}

	creditTransaction := domain.Transaction{
		AccountID:           dofAccount.ID,
		SofNumber:           myAccount.Account,
		DofNumber:           dofAccount.Account,
		TransactionType:     "C",
		Amount:              reqInq.Amount,
		TransactionDateTime: time.Now(),
	}

	err = s.transactionRepository.Insert(ctx, &creditTransaction)

	myAccount.Balance -= reqInq.Amount
	err = s.accountRepository.Update(ctx, &myAccount)
	if err != nil {
		return err
	}

	dofAccount.Balance += reqInq.Amount
	err = s.accountRepository.Update(ctx, &dofAccount)
	if err != nil {
		return err
	}

	// Running goroutines for notification after transfer
	go s.notificationAfterTransfer(myAccount, dofAccount, reqInq.Amount)

	return nil
}

func (s service) notificationAfterTransfer(sofAccount domain.Account, dofAccount domain.Account, amount float64) {
	notificationSender := domain.Notification{
		UserID:    sofAccount.UserID,
		Title:     "Transfer Berhasil",
		Body:      fmt.Sprintf("Transfer senilai %.2f", amount),
		IsRead:    0,
		Status:    1,
		CreatedAt: time.Now(),
	}

	notificationReceiver := domain.Notification{
		UserID:    dofAccount.UserID,
		Title:     "Dana diterima",
		Body:      fmt.Sprintf("Dana diterima senilai %.2f", amount),
		IsRead:    0,
		Status:    1,
		CreatedAt: time.Now(),
	}

	// Insert to DB
	_ = s.notificationRepository.Insert(context.Background(), &notificationSender)
	if channel, ok := s.hub.NotificationChannel[sofAccount.ID]; ok {
		channel <- dto.NotificationData{
			ID:        notificationSender.ID,
			Title:     notificationSender.Title,
			Body:      notificationSender.Body,
			IsRead:    notificationSender.IsRead,
			Status:    notificationSender.Status,
			CreatedAt: notificationSender.CreatedAt,
		}
	}

	_ = s.notificationRepository.Insert(context.Background(), &notificationReceiver)
	if channel, ok := s.hub.NotificationChannel[dofAccount.ID]; ok {
		channel <- dto.NotificationData{
			ID:        notificationReceiver.ID,
			Title:     notificationReceiver.Title,
			Body:      notificationReceiver.Body,
			IsRead:    notificationReceiver.IsRead,
			Status:    notificationReceiver.Status,
			CreatedAt: notificationReceiver.CreatedAt,
		}
	}
}
