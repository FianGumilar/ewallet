package migration

import (
	"log"

	"fiangumilar.id/e-wallet/domain"
	"fiangumilar.id/e-wallet/internal/component"
)

func Migration() {
	err := component.DB.AutoMigrate(
		&domain.User{},
	)
	if err != nil {
		log.Printf("Failed to Migrate: %s", err)
	}
}
