package user

import (
	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
	"github.com/gofiber/fiber/v2"
)

type apiUserAuth struct {
	userService domain.UserService
}

func NewAuth(app *fiber.App, userService domain.UserService, authMid fiber.Handler) {
	api := &apiUserAuth{
		userService: userService,
	}

	app.Post("token/generate", api.GenerateToken)
	app.Get("/token/validate", authMid, api.ValidateToken)
}

func (a apiUserAuth) GenerateToken(ctx *fiber.Ctx) error {
	var req dto.UserReq

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON("Failed to parse request")
	}

	token, err := a.userService.Authenticate(ctx.Context(), req)
	if err != nil {
		return ctx.Status(401).JSON("Failed to authenticate")
	}
	return ctx.Status(200).JSON(token)
}

func (a apiUserAuth) ValidateToken(ctx *fiber.Ctx) error {
	user := ctx.Locals("x-user")

	return ctx.Status(200).JSON(user)
}
