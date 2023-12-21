package migration

import (
	"flag"
	"fmt"
	"strings"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

// Logger is the logger used for the migration package
var Logger = log.New("migration")

var transaction = flag.Bool("trans", false, "Execute the migration inside a single transaction")
var migrateDown = flag.Bool("down", false, "Undo the last migration or undo til the specific --version")
var migrateVersion = flag.String("version", "", "Exec the migrations up/down to the given migration that matches")

// DefaultMigrationOptions contains default options for the gormigrate package
var DefaultMigrationOptions = &gormigrate.Options{
	TableName:      "migrations",
	IDColumnName:   "id",
	IDColumnSize:   255,
	UseTransaction: false,
}

// Run executes the migrations given
func Run(db *gorm.DB, migrations []*gormigrate.Migration) {
	Logger.SetHeader("${time_rfc3339_nano} - [${level}]")
	Logger.DisableColor()
	parseFlags()

	DefaultMigrationOptions.UseTransaction = *transaction

	m := gormigrate.New(db, DefaultMigrationOptions, migrations)

	var err error
	successMsg := "Migrated successfully to version: "
	lastversion := migrations[len(migrations)-1].ID
	if *migrateDown {
		successMsg = "Rolled back to version: "
		if *migrateVersion == "" {
			err = m.RollbackLast()
			lastversion = migrations[len(migrations)-2].ID
		} else {
			err = m.RollbackTo(*migrateVersion)
			lastversion = *migrateVersion
		}
	} else {
		if *migrateVersion == "" {
			err = m.Migrate()
		} else {
			err = m.MigrateTo(*migrateVersion)
			lastversion = *migrateVersion
		}
	}

	if err != nil {
		Logger.Fatalf("Migration failed: %v", err)
	}

	Logger.Info(successMsg + lastversion)
}

// ExecMultiple executes multiple SQL sentences
func ExecMultiple(tx *gorm.DB, sqls string) error {
	for _, sql := range strings.Split(sqls, ";") {
		sql := strings.TrimSpace(sql)
		if sql == "" {
			continue
		}
		if err := tx.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}

func parseFlags() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: go run cmd/migration/main.go [-down] [-version 200601021504]\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()
}
