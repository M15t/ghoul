package user

import (
	"github.com/M15t/ghoul/internal/model"
	dbutil "github.com/M15t/ghoul/pkg/util/db"
	"gorm.io/gorm"
)

// NewDB returns a new user database instance
func NewDB() *DB {
	return &DB{dbutil.NewDB(model.User{})}
}

// DB represents the client for user table
type DB struct {
	*dbutil.DB
}

// FindByUsername queries for single user by username
func (d *DB) FindByUsername(db *gorm.DB, uname string) (*model.User, error) {
	rec := new(model.User)
	if err := d.View(db, rec, "username = ?", uname); err != nil {
		return nil, err
	}
	return rec, nil
}

// FindByRefreshToken queries for single user by refresh token
func (d *DB) FindByRefreshToken(db *gorm.DB, token string) (*model.User, error) {
	rec := new(model.User)
	if err := d.View(db, rec, "refresh_token = ?", token); err != nil {
		return nil, err
	}
	return rec, nil
}
