package auth

import (
	"time"

	"ghoul/internal/model"
	dbutil "ghoul/pkg/util/db"

	"gorm.io/gorm"
)

// New creates new auth service
func New(db *gorm.DB, udb UserDB, jwt JWT, cr Crypter) *Auth {
	return &Auth{
		db:  db,
		udb: udb,
		jwt: jwt,
		cr:  cr,
	}
}

// Auth represents auth application service
type Auth struct {
	db  *gorm.DB
	udb UserDB
	jwt JWT
	cr  Crypter
}

// UserDB represents user repository interface
type UserDB interface {
	dbutil.Intf
	FindByUsername(*gorm.DB, string) (*model.User, error)
	FindByRefreshToken(*gorm.DB, string) (*model.User, error)
}

// JWT represents token generator (jwt) interface
type JWT interface {
	GenerateToken(map[string]interface{}, *time.Time) (string, int, error)
}

// Crypter represents security interface
type Crypter interface {
	CompareHashAndPassword(string, string) bool
	UID() string
}
