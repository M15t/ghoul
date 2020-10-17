package dbutil

import (
	"github.com/jinzhu/gorm"
)

// New creates new database connection to the database server
func New(dialect, dbPsn string, enableLog bool) (*gorm.DB, error) {
	db, err := gorm.Open(dialect, dbPsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Raw("SELECT 1").Rows()
	if err != nil {
		return nil, err
	}

	db.LogMode(enableLog == true)

	return db, nil
}
