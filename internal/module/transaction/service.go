package transaction

import (
	"context"
	"encoding/json"
	"time"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
	"fiangumilar.id/e-wallet/internal/utils"
)

type service struct {
	accountRepository     domain.AccountRepository
	transactionRepository domain.TransactionRepository
	cacheRepository       domain.CacheRepository
}

func NewTransactionService(
	accountRepository domain.AccountRepository,
	transactionRepository domain.TransactionRepository,
	cacheRepository domain.CacheRepository) domain.TransactionService {
	return &service{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		cacheRepository:       cacheRepository,
	}
}

// TransaferInquiry implements domain.TransactionService.
func (s service) TransaferInquiry(ctx context.Context, req dto.TransferInquiryReq) (dto.TransferInquiryRes, error) {
	user := ctx.Value("x-user").(dto.UserData)

	myAccount, err := s.accountRepository.FindByUserID(ctx, user.ID)
	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if myAccount == (domain.Account{}) {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound
	}

	dofAccount, err := s.accountRepository.FindByAccountNumber(ctx, req.AccountNumber)
	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if dofAccount == (domain.Account{}) {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound
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

	dofAccount, err := s.accountRepository.FindByAccountNumber(ctx, reqInq.AccountNumber)
	if err != nil {
		return err
	}

	debitTransaction := domain.Transaction{
		AccountID:           myAccount.ID,
		SofNumber:           myAccount.AccountNumber,
		DofNumber:           dofAccount.AccountNumber,
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
		SofNumber:           myAccount.AccountNumber,
		DofNumber:           dofAccount.AccountNumber,
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

	return nil
}
