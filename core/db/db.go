package db

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"software_updater/core/config"
	"software_updater/core/db/po"
)

var db *gorm.DB

func InitDB(conf *config.DatabaseConfig) (err error) {
	log.Printf("initializing database: %v", conf)

	switch conf.Driver {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(conf.DSN), &gorm.Config{})
	default:
		err = fmt.Errorf("database driver not supported: %s", conf.Driver)
	}
	if err != nil {
		return
	}

	migrator := db.Migrator()
	if !migrator.HasTable(&po.Homepage{}) {
		err = migrator.CreateTable(&po.Homepage{})
		if err != nil {
			return
		}
	}
	if !migrator.HasTable(&po.Version{}) {
		err = migrator.CreateTable(&po.Version{})
		if err != nil {
			return
		}
	}
	if !migrator.HasTable(&po.CurrentVersion{}) {
		err = migrator.CreateTable(&po.CurrentVersion{})
		if err != nil {
			return
		}
	}
	return nil
}
