package dbutil

import (
	"log/slog"

	dbutil "github.com/M15t/ghoul/pkg/util/db"

	"github.com/imdatngo/gowhere"
	sloggorm "github.com/imdatngo/slog-gorm"
	_ "gorm.io/driver/mysql" // DB adapter
	"gorm.io/gorm"
	// EnablePostgreSQL: remove the mysql package above, uncomment the following
	// _ "gorm.io/driver/postgres" // DB adapter
)

// New creates new database connection to the database server
func New(dbPsn string, slogger *slog.Logger) (*gorm.DB, error) {
	// Add your DB related stuffs here, such as:
	// - gorm.DefaultTableNameHandler
	// - gowhere.DefaultConfig
	gowhere.DefaultConfig.Dialect = gowhere.DialectMySQL
	config := new(gorm.Config)

	// Create new slog-gorm instance with slog.Default()
	slogger = slogger.WithGroup("db")
	gConfig := sloggorm.NewConfig(slogger.Handler()).WithTraceAll(true).WithContextKeys(map[string]string{"id": "X-Request-ID"})
	config.Logger = sloggorm.NewWithConfig(gConfig)

	return dbutil.New("mysql", dbPsn, config)

	// EnablePostgreSQL: remove 2 lines above, uncomment the following
	// return dbutil.New("postgres", dbPsn, enableLog)
}

// NewDB creates new DB instance
func NewDB(model interface{}) *dbutil.DB {
	return &dbutil.DB{Model: model}
}
