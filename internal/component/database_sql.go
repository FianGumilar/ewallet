package component

import (
	"database/sql"
	"fmt"
	"log"

	"fiangumilar.id/e-wallet/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

var DbSql *sql.DB

func GetSqlDbConnection(conf *config.Config) *sql.DB {
	var err error

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=True",
		conf.Database.User,
		conf.Database.Pass,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name,
	)

	DbSql, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Failed connect to database: %s", err)
	}
	return DbSql
}
