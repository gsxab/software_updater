package action

import (
	"context"
	"encoding/base64"
	"github.com/tebeka/selenium"
	"io/fs"
	"os"
	"path"
	"software_updater/core/config"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"sync"
	"time"
)

type AppendScreenshot struct {
	Default
	DefaultFactory[AppendScreenshot, *AppendScreenshot]
}

func (a *AppendScreenshot) Path() Path {
	return Path{"browser", "reader", "append_screenshot"}
}

func (a *AppendScreenshot) OutStrNum() int {
	return OneMore
}

func (a *AppendScreenshot) getFilename(name string) string {
	encodedName := base64.URLEncoding.EncodeToString([]byte(name))
	dateSuffix := time.Now().Format("2006-01-02")
	return encodedName + "@" + dateSuffix + ".png"
}

func (a *AppendScreenshot) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
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
	output = AnotherStringToArgs(filename, input)
	return
}

func (a *AppendScreenshot) ToDTO() *DTO {
	return &DTO{
		Output: []string{"filename"},
	}
}
