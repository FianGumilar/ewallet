package notification

import (
	"context"
	"time"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/dto"
	"github.com/gofiber/fiber/v2"
)

type apiNotification struct {
	notificationService domain.NotificationService
}

func NewNotification(app *fiber.App, authMid fiber.Handler, notificationService domain.NotificationService) {
	api := &apiNotification{
		notificationService: notificationService,
	}
	app.Get("/notifications", authMid, api.GetNotification)
}

func (a apiNotification) GetNotification(ctx *fiber.Ctx) error {
	// defer cancel
	c, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := ctx.Locals("x-user").(dto.UserData)

	notification, err := a.notificationService.FindByUser(c, user.ID)
	if err != nil {
		return ctx.SendStatus(401)
	}
	return ctx.Status(200).JSON(notification)
}
