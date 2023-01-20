package action

import (
	"context"
	"github.com/tebeka/selenium"
	"os"
	"path"
	"software_updater/core/config"
	"software_updater/core/db/po"
	"software_updater/core/tools/web"
	"software_updater/core/util/url_util"
	"sync"
)

type CURLSave struct {
	Default
}

func (a *CURLSave) Path() Path {
	return Path{"curl", "access", "curl_save"}
}

func (a *CURLSave) InStrNum() int {
	return 2
}

func (a *CURLSave) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	output = input
	relURL := input.Strings[0]

	baseURL, err := driver.CurrentURL()
	if err != nil {
		return
	}
	url, err := url_util.RelativeURL(baseURL, relURL)
	selCookies, err := driver.GetCookies()
	if err != nil {
		return
	}

	filepath := path.Join(config.Current().Files.CURLSaveDir, input.Strings[1])
	file, err := os.Open(filepath)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	err = web.CURL(url, selCookies, file)
	if err != nil {
		return
	}

	return
}

func (a *CURLSave) ToDTO() *DTO {
	return &DTO{
		OpenPage: true,
		Input:    []string{"rel_url", "file_path"},
	}
}

func (a *CURLSave) NewAction(string) (Action, error) {
	return a, nil
}
