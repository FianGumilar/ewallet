package main

import (
	"log"

	"fiangumilar.id/e-wallet/internal/component"
	"fiangumilar.id/e-wallet/internal/component/cache"
	"fiangumilar.id/e-wallet/internal/component/migration"
	"fiangumilar.id/e-wallet/internal/config"
	"fiangumilar.id/e-wallet/internal/middleware"
	"fiangumilar.id/e-wallet/internal/module/account"
	"fiangumilar.id/e-wallet/internal/module/notification"
	"fiangumilar.id/e-wallet/internal/module/transaction"
	"fiangumilar.id/e-wallet/internal/module/user"
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

	userRepository := user.NewUserRepository(dbConnection)
	accountRepository := account.NewRepository(dbSqlConnection)
	transactionRepository := transaction.NewTransactionRepository(dbConnection)
	notificationRepository := notification.NewRepository(dbSqlConnection)

	userService := user.NewUserService(userRepository, cacheConnection)
	transactionService := transaction.NewTransactionService(accountRepository, transactionRepository, cacheConnection, notificationRepository)
	notificationService := notification.NewNotificationService(notificationRepository)

	authMid := middleware.Authenticate(userService)

	app := fiber.New()

	user.NewAuth(app, userService, authMid)
	transaction.NewTransfer(app, authMid, transactionService)
	notification.NewNotification(app, authMid, notificationService)

	app.Listen(conf.Server.Host + ":" + conf.Server.Port)
}
