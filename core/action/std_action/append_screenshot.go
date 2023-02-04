package std_action

import (
	"context"
	"encoding/base64"
	"github.com/tebeka/selenium"
	"io/fs"
	"os"
	"path"
	"software_updater/core/action"
	"software_updater/core/config"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
	"time"
)

type AppendScreenshot struct {
	action.Default
	action.DefaultFactory[AppendScreenshot, *AppendScreenshot]
}

func (a *AppendScreenshot) Path() action.Path {
	return action.Path{"browser", "reader", "append_screenshot"}
}

func (a *AppendScreenshot) OutStrNum() int {
	return action.OneMore
}

func (a *AppendScreenshot) getFilename(name string) string {
	encodedName := base64.URLEncoding.EncodeToString([]byte(name))
	dateSuffix := time.Now().Format("2006-01-02")
	return encodedName + "@" + dateSuffix + ".png"
}

func (a *AppendScreenshot) Do(ctx context.Context, driver selenium.WebDriver, input *action.Args, version *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	bytes, err := driver.Screenshot()
	if err != nil {
		logs.Error(ctx, "selenium screenshot failed", err)
		return
	}
	filename := a.getFilename(version.Name)
	pathname := path.Join(config.Current().Files.ScreenshotDir, filename)
	err = os.WriteFile(pathname, bytes, fs.FileMode(0o644))
	if err != nil {
		logs.Error(ctx, "screenshot saving failed", err, "filename", filename, "pathname", pathname)
		return
	}
	output = action.AnotherStringToArgs(filename, input)
	return
}

func (a *AppendScreenshot) ToDTO() *action.DTO {
	return &action.DTO{
		Output: []string{"filename"},
	}
}
