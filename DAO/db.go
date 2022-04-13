package DAO

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var db *gorm.DB

func Init() {

	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.", k)
		}
		return v
	}

	// Local environment setup: https://github.com/GoogleCloudPlatform/golang-samples/tree/main/cloudsql/mysql/database-sql
	var (
		dbUser         = mustGetenv("DB_USER")              // e.g. 'my-db-user'
		dbPwd          = mustGetenv("DB_PASS")              // e.g. 'my-db-password'
		dbName         = mustGetenv("DB_NAME")              // e.g. 'my-database'
		unixSocketPath = mustGetenv("INSTANCE_UNIX_SOCKET") // e.g. '/cloudsql/project:region:instance'
	)

	dsn := fmt.Sprintf("%s:%s@unix(%s)/%s?parseTime=true",
		dbUser, dbPwd, unixSocketPath, dbName)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(fmt.Errorf("failed to connect database %v", err))
	}

	// Migrate the schema
	db.AutoMigrate(
		&Course{},
		&CourseGrade{},
		&Question{},
		&Answer{},
	)

}

func DB() *gorm.DB {
	return db
}
