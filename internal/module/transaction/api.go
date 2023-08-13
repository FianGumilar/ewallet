package transaction

import (
	"log"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
	"github.com/gofiber/fiber/v2"
)

type apiTransaction struct {
	transactionService domain.TransactionService
}

func NewTransfer(app *fiber.App, authMid fiber.Handler, transactionService domain.TransactionService) {
	api := apiTransaction{
		transactionService: transactionService,
	}

	app.Post("transfer/inquiry", authMid, api.TransferInquiry)
	app.Post("transfer/execute", authMid, api.TransferExecute)
}

func (a apiTransaction) TransferInquiry(ctx *fiber.Ctx) error {
	var req dto.TransferInquiryReq

	if err := ctx.BodyParser(&req); err != nil {
		log.Printf("error %s", err)
		return ctx.Status(400).JSON(dto.Response{
			Message: "Invalid body Request",
		})
	}

	inquiry, err := a.transactionService.TransferInquiry(ctx.Context(), req)
	if err != nil {
		log.Printf("error %s", err)
		return ctx.Status(400).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return ctx.Status(200).JSON(inquiry)
}

func (a apiTransaction) TransferExecute(ctx *fiber.Ctx) error {
	var req dto.TransferExecuteReq

	if err := ctx.BodyParser(&req); err != nil {
		log.Printf("error %s", err)
		return ctx.Status(400).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	err := a.transactionService.TransferExecute(ctx.Context(), req)
	if err != nil {
		log.Printf("error %s", err)
		return ctx.Status(400).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return ctx.SendStatus(200)
}
