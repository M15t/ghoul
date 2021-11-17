package country

import (
	"ghoul/internal/model"
	"ghoul/pkg/rbac"
	dbutil "ghoul/pkg/util/db"

	"gorm.io/gorm"
)

// New creates new country application service
func New(db *gorm.DB, cdb dbutil.Intf, rbacSvc rbac.Intf) *Country {
	return &Country{
		db:   db,
		cdb:  cdb,
		rbac: rbacSvc,
	}
}

// Country represents country application service
type Country struct {
	db   *gorm.DB
	cdb  dbutil.Intf
	rbac rbac.Intf
}

// NewDB returns a new country database instance
func NewDB() *dbutil.DB {
	return dbutil.NewDB(model.Country{})
}
