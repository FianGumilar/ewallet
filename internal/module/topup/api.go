package topup

import (
	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
	"github.com/gofiber/fiber/v2"
)

type apiTopUp struct {
	topUpService domain.TopUpService
}

func NewTopUpApi(app *fiber.App, authMid fiber.Handler, topUpService domain.TopUpService) {
	a := &apiTopUp{
		topUpService: topUpService,
	}

	app.Post("topup/initialize", authMid, a.InitializeTopUp)
}

func (a apiTopUp) InitializeTopUp(ctx *fiber.Ctx) error {
	var req dto.TopUpRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	user := ctx.Locals("x-user").(dto.UserData)
	req.UserID = user.ID

	res, err := a.topUpService.InitializeTopUp(ctx.Context(), req)
	if err != nil {
		return ctx.Status(401).JSON(dto.Response{
			Message: err.Error(),
		})
	}
	return ctx.Status(200).JSON(res)
}
