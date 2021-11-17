package dbutil

import (
	dbutil "ghoul/pkg/util/db"

	"github.com/imdatngo/gowhere"
	_ "gorm.io/driver/mysql" // DB adapter
	"gorm.io/gorm"
	// EnablePostgreSQL: remove the mysql package above, uncomment the following
	// _ "gorm.io/driver/postgres" // DB adapter
)

// New creates new database connection to the database server
func New(dbPsn string, enableLog bool) (*gorm.DB, error) {
	// Add your DB related stuffs here, such as:
	// - gorm.DefaultTableNameHandler
	// - gowhere.DefaultConfig
	gowhere.DefaultConfig.Dialect = gowhere.DialectMySQL
	return dbutil.New("mysql", dbPsn, enableLog)

	// EnablePostgreSQL: remove 2 lines above, uncomment the following
	// return dbutil.New("postgres", dbPsn, enableLog)
}

// NewDB creates new DB instance
func NewDB(model interface{}) *dbutil.DB {
	return &dbutil.DB{Model: model}
}
