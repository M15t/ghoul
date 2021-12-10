package user

import (
	"github.com/M15t/ghoul/internal/model"
	"github.com/M15t/ghoul/pkg/rbac"
	dbutil "github.com/M15t/ghoul/pkg/util/db"

	"gorm.io/gorm"
)

// New creates new user application service
func New(db *gorm.DB, udb MyDB, rbacSvc rbac.Intf, cr Crypter) *User {
	return &User{db: db, udb: udb, rbac: rbacSvc, cr: cr}
}

// User represents user application service
type User struct {
	db   *gorm.DB
	udb  MyDB
	rbac rbac.Intf
	cr   Crypter
}

// MyDB represents user repository interface
type MyDB interface {
	dbutil.Intf
	FindByUsername(*gorm.DB, string) (*model.User, error)
}

// Crypter represents security interface
type Crypter interface {
	CompareHashAndPassword(hasedPwd string, rawPwd string) bool
	HashPassword(string) string
}
