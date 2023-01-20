package main

import (
	"context"
	"log"
	"software_updater/core/config"
	"software_updater/core/db"
	"software_updater/core/engine"
	"software_updater/core/tools/web"
)

func main() {
	conf, err := config.Load("./conf.yaml")
	if err != nil {
		log.Panic(err)
	}

	err = db.InitDB(&conf.Database)
	if err != nil {
		log.Panic(err)
	}

	err = web.InitSelenium(&conf.Selenium)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		_ = web.StopSelenium()
	}()

	e, err := engine.InitEngine(&conf.Engine)
	if err != nil {
		log.Panic(err)
	}

	err = e.CrawlAll(context.Background())
	if err != nil {
		log.Panic(err)
	}
}
