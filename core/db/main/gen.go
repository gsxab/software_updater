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
	g.GenerateAllTable()
	g.ApplyBasic(po.Homepage{}, po.Version{}, po.CurrentVersion{})
	g.Execute()
}
