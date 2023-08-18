package midtrans

import (
	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
	"github.com/gofiber/fiber/v2"
)

type apiMidtrans struct {
	midtransService domain.MidtransService
	topUpService    domain.TopUpService
}

func NewMidtransApi(app *fiber.App, midtransService domain.MidtransService, topUpService domain.TopUpService) {
	api := &apiMidtrans{
		topUpService:    topUpService,
		midtransService: midtransService,
	}

	app.Post("midtrans/payment-callback", api.paymentHandlerNotification)
}

func (a apiMidtrans) paymentHandlerNotification(ctx *fiber.Ctx) error {
	var notificationPayload map[string]interface{}

	if err := ctx.BodyParser(&notificationPayload); err != nil {
		return ctx.Status(400).JSON(dto.Response{Message: err.Error()})
	}

	orderId, exist := notificationPayload["order_id"].(string)
	if !exist {
		return ctx.SendStatus(400)
	}

	success, _ := a.midtransService.VerifyPayment(ctx.Context(), notificationPayload)
	if success {
		_ = a.topUpService.ConfirmedTopUp(ctx.Context(), orderId)
		return ctx.SendStatus(200)
	}
	return ctx.SendStatus(400)
}
