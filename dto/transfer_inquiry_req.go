package dto

type TransferInquiryReq struct {
	Account string  `json:"account"`
	Amount  float64 `json:"amount"`
}
