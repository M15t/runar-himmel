package main

import (
	"fmt"
	"log"
	"runar-himmel/config"
	"runar-himmel/internal/db"
	"runar-himmel/internal/types"
	"runar-himmel/pkg/util/crypter"
	"runar-himmel/pkg/util/migration"
	"time"

	"runar-himmel/internal/rbac"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// EnablePostgreSQL: remove this and all tx.Set() functions bellow
// var defaultTableOpts = "ENGINE=InnoDB ROW_FORMAT=DYNAMIC"
var defaultTableOpts = ""

func main() {
	if config.IsLambda() {
		// start lambda request handler
		lambda.Start(handler)
		return
	}

	// start the function directly
	if err := Run(); err != nil {
		log.Println(err)
	}
}

func handler() (string, error) {
	if err := Run(); err != nil {
		return "DB Migration failed!", err
	}
	return "DB Migration completed!", nil
}

// Run executes the migration
func Run() (respErr error) {
	cfg, err := config.LoadAll()
	if err != nil {
		return err
	}

	db, sqldb, err := db.New(cfg.DB)
	if err != nil {
		return err
	}
	defer sqldb.Close()
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
				respErr = fmt.Errorf("unknown error: %+v", x)
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
		// create initial table(s)
		{
			ID: "202312151405",
			Migrate: func(tx *gorm.DB) error {
				// it's a good pratice to copy the struct inside the function,
				// so side effects are prevented if the original struct changes during the time

				// Drop existing table if there is, in case re-apply this migration
				if err := tx.Migrator().DropTable(&types.User{}); err != nil {
					return err
				}

				if err := tx.Set("gorm:table_options", defaultTableOpts).AutoMigrate(&types.User{}); err != nil {
					return err
				}

				// insert default users
				now := time.Now()
				defaultUsers := []*types.User{
					{
						Email:           "odin@runar-himmel.sky",
						EmailVerifiedAt: &now,
						Phone:           "+6281234567890",
						PhoneVerifiedAt: &now,
						FirstName:       "Odin",
						LastName:        "Allfather",
						Role:            rbac.RoleSuperAdmin,
					},
					{
						Email:           "thor@runar-himmel.sky",
						EmailVerifiedAt: &now,
						Phone:           "+6281234567891",
						PhoneVerifiedAt: &now,
						FirstName:       "Thor",
						LastName:        "God of Thunder",
						Role:            rbac.RoleAdmin,
					},
					{
						Email:           "loki@runar-himmel.sky",
						EmailVerifiedAt: &now,
						Phone:           "+6281234567892",
						PhoneVerifiedAt: &now,
						FirstName:       "Loki",
						LastName:        "Laufeyjarson",
						Role:            rbac.RoleUser,
					},
				}
				for _, usr := range defaultUsers {
					if usr.Password == "" {
						usr.Password = usr.Role + "123!@#"
					}
					usr.Password = crypter.HashPassword(usr.Password)
					if err := tx.Create(usr).Error; err != nil {
						return err
					}
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users")
			},
		},
		// create "sessions" table
		{
			ID: "202312211623",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.Set("gorm:table_options", defaultTableOpts).AutoMigrate(&types.Session{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("sessions")
			},
		},
	})

	return nil
}
