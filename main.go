package main

import (
	"context"
	"log"
	"software_updater/core/config"
	"software_updater/core/db"
	"software_updater/core/engine"
	"software_updater/core/logs"
	"software_updater/core/tools/web"
	"software_updater/ui/client"
	"software_updater/ui/webui"
)

func main() {
	conf, err := config.Load("./conf.yaml")
	if err != nil {
		log.Panic(err)
	}

	err = db.InitDB(conf.Database)
	if err != nil {
		log.Panic(err)
	}

	err = web.InitSelenium(conf.Selenium)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		_ = web.StopSelenium()
	}()

	e, err := engine.InitEngine(conf.Engine)
	if err != nil {
		log.Panic(err)
	}

	uiMode := conf.Extra["ui_mode"]
	switch uiMode {
	case "client":
		logs.WarnM(context.Background(), "client ui selected, which is not fully usable in current version!!!")
		err = client.InitAndRun(context.Background(), conf.Extra["client_ui_setting"])
		if err != nil {
			log.Panic(err)
		}
	case "", "web":
		logs.InfoM(context.Background(), "web ui selected")
		err = webui.InitAndRun(context.Background(), conf.Extra["web_ui_setting"])
		if err != nil {
			log.Panic(err)
		}
	case "off":
		err = e.RunAll(context.Background())
		if err != nil {
			log.Panic(err)
		}
	}
}
