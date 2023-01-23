package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/util/version_util"
	"sync"
)

type CheckLaterVersion struct {
	Default
	DefaultFactory[CheckLaterVersion, *CheckLaterVersion]
	VersionFormat string `json:"format"`
}

func (a *CheckLaterVersion) Path() Path {
	return Path{"basic", "value_check", "version_gt"}
}

func (a *CheckLaterVersion) Do(_ context.Context, _ selenium.WebDriver, _ *Args, version *po.Version, wg *sync.WaitGroup) (output *Args, exit Result, err error) {
	if version.Previous == nil {
		return
	}
	previousVersion, err := version_util.Parse(a.VersionFormat, *version.Previous)
	if err != nil {
		return
	}
	currentVersion, err := version_util.Parse(a.VersionFormat, version.Version)
	if err != nil {
		return
	}
	if previousVersion.LT(currentVersion) {
		exit = StopFlow
	}
	return
}

func (a *CheckLaterVersion) ToDTO() *DTO {
	return &DTO{
		Values: map[string]string{"format": a.VersionFormat},
	}
}
