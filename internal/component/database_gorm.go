package component

import (
	"fmt"
	"log"

	"fiangumilar.id/e-wallet/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetConnectionDB(conf *config.Config) *gorm.DB {
	var err error

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=True",
		conf.Database.User,
		conf.Database.Pass,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name,
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed Connect to Database %s", err)
	}
	return DB
}
