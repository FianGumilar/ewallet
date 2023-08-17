package midtrans

import (
	"context"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/internal/config"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type service struct {
	client         snap.Client
	midtransConfig config.Midtrans
}

func NewMidtransService(conf *config.Config) domain.MidtransService {
	var client snap.Client

	envi := midtrans.Sandbox
	if conf.Midtrans.IsProd {
		envi = midtrans.Production
	}

	client.New(conf.Midtrans.Key, envi)

	return &service{
		client:         client,
		midtransConfig: conf.Midtrans,
	}
}

// GenerateSnapURL implements domain.MidtransService.
func (s service) GenerateSnapURL(ctx context.Context, t *domain.TopUp) error {
	// Initiate Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  t.ID,
			GrossAmt: int64(t.Amount),
		},
	}

	// Request create Snap transaction to Midtrans
	snapResp, err := s.client.CreateTransaction(req)
	if err != nil {
		return err
	}

	t.SnapURL = snapResp.RedirectURL
	return nil
}

// VerifyPayment implements domain.MidtransService.
func (s service) VerifyPayment(ctx context.Context, data map[string]interface{}) (bool, error) {
	var client coreapi.Client
	envi := midtrans.Sandbox
	if s.midtransConfig.IsProd {
		envi = midtrans.Production
	}

	client.New(s.midtransConfig.Key, envi)

	// Get order-id from payload
	orderId, exists := data["order_id"].(string)
	if !exists {
		// do something when key `order_id` not found
		return false, domain.ErrInvalidPayload
	}

	// Check transaction to Midtrans with param orderId
	transactionStatusResp, e := client.CheckTransaction(orderId)
	if e != nil {
		return false, e
	} else {
		if transactionStatusResp != nil {
			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				} else if transactionStatusResp.FraudStatus == "accept" {
					// TODO set transaction status on your database to 'success'
					return true, nil
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				// TODO set transaction status on your databaase to 'success'
				return true, nil
			} else if transactionStatusResp.TransactionStatus == "deny" {
				// TODO you can ignore 'deny', because most of the time it allows payment retries
				// and later can become success
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				// TODO set transaction status on your databaase to 'failure'
			} else if transactionStatusResp.TransactionStatus == "pending" {
				// TODO set transaction status on your databaase to 'pending' / waiting payment
			}
		}
	}
	return false, nil
}
