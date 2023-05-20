package dbutil

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// New creates new database connection to the database server
func New(dialect, dbPsn string, cfg *gorm.Config) (db *gorm.DB, err error) {
	switch dialect {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dbPsn), cfg)
		if err != nil {
			return nil, err
		}
	case "postgres":
		db, err = gorm.Open(postgres.Open(dbPsn), cfg)
		if err != nil {
			return nil, err
		}
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(dbPsn), cfg)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
