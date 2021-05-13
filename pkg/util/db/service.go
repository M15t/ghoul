package dbutil

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New creates new database connection to the database server
func New(dialect, dbPsn string, enableLog bool) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	db := new(gorm.DB)
	switch dialect {
	case "mysql":
		db, err := gorm.Open(mysql.Open(dbPsn), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			return nil, err
		}
		return db, nil
	case "postgres":
		db, err := gorm.Open(postgres.Open(dbPsn), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	return db, nil
}
