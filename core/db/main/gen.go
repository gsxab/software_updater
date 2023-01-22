package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
	"software_updater/core/db/po"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./core/db/dao",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})
	db, err := gorm.Open(sqlite.Open("./software.db"))
	if err != nil {
		panic(err)
	}
	g.UseDB(db)
	g.ApplyBasic(po.Homepage{}, po.Version{}, po.CurrentVersion{})
	g.Execute()

	if db.Migrator().HasTable(&po.Homepage{}) {
		err = db.Migrator().AutoMigrate(po.Homepage{})
	} else {
		err = db.Migrator().CreateTable(&po.Homepage{})
	}
	if err != nil {
		panic(err)
	}

	if db.Migrator().HasTable(&po.Version{}) {
		err = db.Migrator().AutoMigrate(&po.Version{})
	} else {
		err = db.Migrator().CreateTable(&po.Version{})
	}
	if err != nil {
		panic(err)
	}

	if db.Migrator().HasTable(&po.CurrentVersion{}) {
		err = db.Migrator().AutoMigrate(&po.CurrentVersion{})
	} else {
		err = db.Migrator().CreateTable(&po.CurrentVersion{})
	}
	if err != nil {
		panic(err)
	}
}
