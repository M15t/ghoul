package dbutil

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New creates new database connection to the database server
func New(dialect, dbPsn string, cfg *gorm.Config) (*gorm.DB, error) {
	db := new(gorm.DB)
	switch dialect {
	case "mysql":
		db, err := gorm.Open(mysql.Open(dbPsn), cfg)
		if err != nil {
			return nil, err
		}
		return db, nil
	case "postgres":
		db, err := gorm.Open(postgres.Open(dbPsn), cfg)
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	return db, nil
}
