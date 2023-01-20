package web

import (
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"software_updater/core/config"
)

var service *selenium.Service
var driver selenium.WebDriver

func InitSelenium(conf *config.SeleniumConfig) (err error) {
	log.Printf("initializing selenium: %v", conf)

	// Run Chrome browser
	service, err = selenium.NewChromeDriverService(conf.DriverPath, 4444)
	if err != nil {
		return err
	}

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: conf.Params})

	driver, err = selenium.NewRemote(caps, "")
	if err != nil {
		return err
	}

	return nil
}

func Driver() selenium.WebDriver {
	return driver
}

func StopSelenium() error {
	err := driver.Quit()
	if err != nil {
		return err
	}
	err = service.Stop()
	return err
}
