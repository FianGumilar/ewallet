package user

import (
	"log"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
	"github.com/gofiber/fiber/v2"
)

type apiUserAuth struct {
	userService domain.UserService
	fdsService  domain.FdsService
}

func NewAuth(app *fiber.App, userService domain.UserService, fdsService domain.FdsService, authMid fiber.Handler) {
	api := &apiUserAuth{
		userService: userService,
		fdsService:  fdsService,
	}

	app.Post("token/generate", api.GenerateToken)
	app.Get("/token/validate", authMid, api.ValidateToken)
	app.Post("/user/register", api.RegisterUser)
	app.Post("/user/validate", api.ValidateOtp)
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

	if !a.fdsService.IsAuthorized(ctx.Context(), ctx.Get("X-FORWARDED-FOR"), token.UserId) {
		return ctx.SendStatus(401)
	}

	return ctx.Status(200).JSON(token)
}

func (a apiUserAuth) ValidateToken(ctx *fiber.Ctx) error {
	user := ctx.Locals("x-user")

	return ctx.Status(200).JSON(user)
}

func (a apiUserAuth) RegisterUser(ctx *fiber.Ctx) error {
	var req dto.UserRegisterReq

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	res, err := a.userService.Register(ctx.Context(), req)
	if err != nil {
		return ctx.SendStatus(400)
	}

	return ctx.Status(200).JSON(res)
}

func (a apiUserAuth) ValidateOtp(ctx *fiber.Ctx) error {
	var req dto.ValidateOtpReq

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON("error parsing request")
	}

	err := a.userService.ValidateOtp(ctx.Context(), req)
	if err != nil {
		log.Printf("error %s", err)
		return ctx.SendStatus(400)
	}
	return ctx.SendStatus(200)
}
