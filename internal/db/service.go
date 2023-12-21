package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/imdatngo/gowhere"

	"runar-himmel/config"
	dbutil "runar-himmel/pkg/util/db"

	// _ "gorm.io/driver/sqlite" // DB adapter
	// _ "gorm.io/gorm/dialects/postgres" // DB adapter
	_ "gorm.io/driver/mysql" // DB adapter
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// EnablePostgreSQL: remove the mysql package above, uncomment the following
)

// New creates new database connection to the database server
func New(cfg config.DB) (db *gorm.DB, sqldb *sql.DB, err error) {
	// Add your DB related stuffs here, such as:
	// - gorm.DefaultTableNameHandler
	// - gowhere.DefaultConfig

	// ! EnablePostgreSQL
	// gowhere.DefaultConfig.Dialect = gowhere.DialectPostgreSQL
	gowhere.DefaultConfig.Dialect = gowhere.DialectMySQL

	// logger config
	var lo logger.Interface
	if cfg.Logging > 0 {
		lo = logger.Default.LogMode(logger.LogLevel(cfg.Logging))
	} else {
		lo = logger.Discard
	}

	// ! EnablePostgreSQL: remove 2 lines above, uncomment the following
	// db, err := dbutil.New("postgres", dbPsn, enableLog)
	db, err = dbutil.New("mysql", cfg, &gorm.Config{
		Logger:                                   lo,
		AllowGlobalUpdate:                        false,
		CreateBatchSize:                          1000,
		DisableForeignKeyConstraintWhenMigrating: true,
		// NamingStrategy:                           schema.NamingStrategy{
		// 	// SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		// },
	})

	// connection pool settings
	sqldb, err = db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot get generic db instance: %w", err)
	}
	//! NOTE: These are not one-size-fits-all settings. Turn it based on your db settings!
	sqldb.SetMaxIdleConns(10)
	sqldb.SetMaxOpenConns(10)
	sqldb.SetConnMaxLifetime(30 * time.Minute)
	sqldb.SetConnMaxIdleTime(10 * time.Minute)

	return
}
