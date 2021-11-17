package migration

import (
	"fmt"
	"os"
	"time"

	"ghoul/config"
	"ghoul/internal/model"
	dbutil "ghoul/internal/util/db"
	"ghoul/pkg/util/crypter"
	"ghoul/pkg/util/migration"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// EnablePostgreSQL: remove this and all tx.Set() functions bellow
var defaultTableOpts = "ENGINE=InnoDB ROW_FORMAT=DYNAMIC"

// Base represents base columns for all tables. Do not use gorm.Model because of uint ID
type Base struct {
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Run executes the migration
func Run() (respErr error) {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := dbutil.New(cfg.DbPsn, false)
	if err != nil {
		return err
	}
	// connection.Close() is not available for GORM 1.20.0
	// defer db.Close()

	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				respErr = fmt.Errorf("%s", x)
			case error:
				respErr = x
			default:
				respErr = fmt.Errorf("Unknown error: %+v", x)
			}
		}
	}()

	// EnablePostgreSQL: remove these
	// workaround for "Index column size too large" error on migrations table
	initSQL := "CREATE TABLE IF NOT EXISTS migrations (id VARCHAR(255) PRIMARY KEY) " + defaultTableOpts
	if err := db.Exec(initSQL).Error; err != nil {
		return err
	}

	migration.Run(db, []*gormigrate.Migration{
		// create initial tables
		{
			ID: "201905051012",
			Migrate: func(tx *gorm.DB) error {
				// it's a good pratice to copy the struct inside the function,
				// so side effects are prevented if the original struct changes during the time

				type Country struct {
					Base
					Name      string `gorm:"type:varchar(255)"`
					Code      string `gorm:"type:varchar(10)"`
					PhoneCode string `gorm:"type:varchar(10)"`
				}

				type User struct {
					Base
					FirstName    string `gorm:"type:varchar(255)"`
					LastName     string `gorm:"type:varchar(255)"`
					Email        string `gorm:"type:varchar(255)"`
					Mobile       string `gorm:"type:varchar(255)"`
					Username     string `gorm:"type:varchar(255);unique_index;not null"`
					Password     string `gorm:"type:varchar(255);not null"`
					LastLogin    *time.Time
					Blocked      bool   `gorm:"not null;default:0"`
					RefreshToken string `gorm:"type:varchar(255)"`
					Role         string `gorm:"type:varchar(255)"`
				}

				if err := tx.Set("gorm:table_options", defaultTableOpts).AutoMigrate(&Country{}, &User{}); err != nil {
					return err
				}

				// insert default users
				defaultUsers := []*User{
					{
						Username: "superadmin",
						Password: os.Getenv("SUPERADMIN_PWD"),
						Email:    "superadmin@ghoul.com",
						Role:     model.RoleSuperAdmin,
					},
					{
						Username: "admin",
						Password: os.Getenv("ADMIN_PWD"),
						Email:    "admin@ghoul.com",
						Role:     model.RoleAdmin,
					},
					{
						Username: "user",
						Password: os.Getenv("USER_PWD"),
						Email:    "user@ghoul.com",
						Role:     model.RoleUser,
					},
				}
				for _, usr := range defaultUsers {
					if usr.Password == "" {
						usr.Password = usr.Username + "123!@#"
					}
					usr.Password = crypter.HashPassword(usr.Password)
					if err := tx.Create(usr).Error; err != nil {
						return err
					}
				}

				// insert default countries
				defaultCountries := []*Country{
					{
						Name:      "Singapore",
						Code:      "SG",
						PhoneCode: "+65",
					},
				}
				for _, rec := range defaultCountries {
					if err := tx.Create(rec).Error; err != nil {
						return err
					}
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users", "countries")
			},
		},
	})

	return nil
}
