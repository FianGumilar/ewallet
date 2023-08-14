package sse

import (
	"bufio"
	"encoding/json"
	"fmt"

	"fiangumilar.id/e-wallet/dto"
	"github.com/gofiber/fiber/v2"
)

type notificationSse struct {
	hub *dto.Hub
}

func NewNotificationSse(app *fiber.App, authMid fiber.Handler, hub *dto.Hub) {
	api := &notificationSse{
		hub: hub,
	}
	app.Get("sse/notification-stream", authMid, api.StreamNotification)
}

func (n notificationSse) StreamNotification(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/event-stream")

	user := ctx.Locals("x-user").(dto.UserData)
	n.hub.NotificationChannel[user.ID] = make(chan dto.NotificationData)

	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		event := fmt.Sprintf("event: %s\n"+"data: \n\n", "initial")

		_, _ = fmt.Fprint(w, event)
		_ = w.Flush()

		// Looping Notification Channel
		for notification := range n.hub.NotificationChannel[user.ID] {
			data, _ := json.Marshal(notification)

			event = fmt.Sprintf("event: %s\n"+"data: %s\n\n", "notification-updated", data)
			_, _ = fmt.Fprint(w, event)
			_ = w.Flush()
		}
	})
	return nil
}
