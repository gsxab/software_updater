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
	"sync"
	"time"
)

type AppendScreenshot struct {
	Default
	DefaultFactory[AppendScreenshot, *AppendScreenshot]
	Name string `json:"name"`
}

func (a *AppendScreenshot) Path() Path {
	return Path{"browser", "reader", "append_screenshot"}
}

func (a *AppendScreenshot) OutStrNum() int {
	return OneMore
}

func (a *AppendScreenshot) getFilename() string {
	encodedName := base64.URLEncoding.EncodeToString([]byte(a.Name))
	dateSuffix := time.Now().Format("2006-01-02")
	return encodedName + "@" + dateSuffix + ".png"
}

func (a *AppendScreenshot) Do(ctx context.Context, driver selenium.WebDriver, input *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	bytes, err := driver.Screenshot()
	if err != nil {
		return
	}
	filename := a.getFilename()
	err = os.WriteFile(path.Join(config.Current().Files.ScreenshotDir, filename), bytes, fs.FileMode(0o644))
	if err != nil {
		return
	}
	output = AnotherStringToArgs(filename, input)
	return
}

func (a *AppendScreenshot) ToDTO() *DTO {
	return &DTO{
		Output: []string{"filename"},
		Values: map[string]string{"name": a.Name},
	}
}
