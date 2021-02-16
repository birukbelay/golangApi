package db

import (
	"os"

	"github.com/jinzhu/gorm"
)

func PgDbConn() *gorm.DB {
	dbconn, err := gorm.Open("postgres", os.ExpandEnv("postgres://${PG_USER_NAME}:${PG_PASS}@localhost/${PG_DB_NAME}?sslmode=disable"))
	// createTables(dbconn)
	if err != nil {
		panic(err)
	}
	return dbconn
}



