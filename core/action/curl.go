package action

import (
	"bytes"
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/tools/web"
	"software_updater/core/util/url_util"
	"sync"
)

type CURL struct {
	Default
}

func (a *CURL) Path() Path {
	return Path{"curl", "access", "curl_content"}
}

func (a *CURL) InStrNum() int {
	return 1
}

func (a *CURL) OutStrNum() int {
	return 1
}

func (a *CURL) Do(ctx context.Context, driver selenium.WebDriver, input *Args, _ *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
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
		logs.Error(ctx, "selenium cookie failed", err)
		return
	}

	buffer := new(bytes.Buffer)
	err = web.CURL(url, selCookies, buffer)
	if err != nil {
		logs.Error(ctx, "cURL failed", err, "URL", url)
		return
	}

	result := buffer.String()
	output = StringToArgs(result, input)
	return
}

func (a *CURL) ToDTO() *DTO {
	return &DTO{
		OpenPage: true,
		Input:    []string{"rel_url"},
		Output:   []string{"remote_data"},
	}
}

func (a *CURL) NewAction(_ string) (Action, error) {
	return a, nil
}
