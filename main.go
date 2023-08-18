package main

import (
	"log"

	"fiangumilar.id/e-wallet/dto"
	"fiangumilar.id/e-wallet/internal/component"
	"fiangumilar.id/e-wallet/internal/component/cache"
	"fiangumilar.id/e-wallet/internal/component/migration"
	"fiangumilar.id/e-wallet/internal/config"
	"fiangumilar.id/e-wallet/internal/middleware"
	"fiangumilar.id/e-wallet/internal/module/account"
	"fiangumilar.id/e-wallet/internal/module/midtrans"
	"fiangumilar.id/e-wallet/internal/module/notification"
	"fiangumilar.id/e-wallet/internal/module/template"
	"fiangumilar.id/e-wallet/internal/module/topup"
	"fiangumilar.id/e-wallet/internal/module/transaction"
	"fiangumilar.id/e-wallet/internal/module/user"
	"fiangumilar.id/e-wallet/internal/sse"
	"github.com/gofiber/fiber/v2"
)

func main() {
	conf := config.Get()

	dbConnection := component.GetConnectionDB(conf)
	dbSqlConnection := component.GetSqlDbConnection(conf)

	//Migration
	migration.Migration()

	cacheConnection, err := cache.NewRedisClient(conf)
	if err != nil {
		log.Printf("Failed connecting to redis: %v", err)
	}
	log.Println("Successfully connected Redis")

	hub := &dto.Hub{
		NotificationChannel: map[int64]chan dto.NotificationData{},
	}

	userRepository := user.NewUserRepository(dbConnection)
	accountRepository := account.NewRepository(dbSqlConnection)
	transactionRepository := transaction.NewTransactionRepository(dbConnection)
	notificationRepository := notification.NewRepository(dbSqlConnection)
	templateRepository := template.NewTemplateRepository(dbSqlConnection)
	topUpRepository := topup.NewTopUpRepository(dbSqlConnection)

	userService := user.NewUserService(userRepository, cacheConnection)
	notificationService := notification.NewNotificationService(notificationRepository, templateRepository, hub)
	transactionService := transaction.NewTransactionService(accountRepository, transactionRepository, cacheConnection, notificationService)
	midtransService := midtrans.NewMidtransService(conf)
	topUpService := topup.NewTopUpService(notificationService, midtransService, topUpRepository, accountRepository)

	authMid := middleware.Authenticate(userService)

	app := fiber.New()

	user.NewAuth(app, userService, authMid)
	transaction.NewTransfer(app, authMid, transactionService)
	notification.NewNotification(app, authMid, notificationService)
	topup.NewTopUpApi(app, authMid, topUpService)
	midtrans.NewMidtransApi(app, midtransService, topUpService)

	sse.NewNotificationSse(app, authMid, hub)

	app.Listen(conf.Server.Host + ":" + conf.Server.Port)
}
