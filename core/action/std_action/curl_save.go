package std_action

import (
	"context"
	"github.com/tebeka/selenium"
	"os"
	"path"
	"software_updater/core/action"
	"software_updater/core/config"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/tools/web"
	"software_updater/core/util/url_util"
	"sync"
)

type CURLSave struct {
	action.Default
}

func (a *CURLSave) Path() action.Path {
	return action.Path{"curl", "access", "curl_save"}
}

func (a *CURLSave) InStrNum() int {
	return 2
}

func (a *CURLSave) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, _ *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	output = input
	relURL := input.Strings[0]

	baseURL, err := driver.CurrentURL()
	if err != nil {
		logs.Error(ctx, "selenium current url failed", err)
		return
	}
	url, err := url_util.RelativeURL(baseURL, relURL)
	if err != nil {
		logs.Error(ctx, "relative url calculation failed", err, "baseURL", baseURL, "relURL", relURL)
		return
	}
	selCookies, err := driver.GetCookies()
	if err != nil {
		logs.Error(ctx, "selenium cookies failed", err)
		return
	}

	pathname := path.Join(config.Current().Files.CURLSaveDir, input.Strings[1])
	file, err := os.Open(pathname)
	if err != nil {
		logs.Error(ctx, "file opening failed", err, "filename", input.Strings[1], "pathname", pathname)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logs.Error(ctx, "close file failed", err)
		}
	}(file)

	err = web.CURL(url, selCookies, file)
	if err != nil {
		logs.Error(ctx, "cURL failed", err, "URL", url)
		return
	}

	return
}

func (a *CURLSave) ToDTO() *action.DTO {
	return &action.DTO{
		OpenPage: true,
		Input:    []string{"rel_url", "file_path"},
	}
}

func (a *CURLSave) NewAction(string) (action.Action, error) {
	return a, nil
}
