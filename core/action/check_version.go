package action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/util/version_util"
	"sync"
)

type CheckVersion struct {
	Default
	DefaultFactory[CheckVersion, *CheckVersion]
	VersionFormat string `json:"format"`
}

func (a *CheckVersion) Path() Path {
	return Path{"basic", "value_check", "version_neq"}
}

func (a *CheckVersion) Do(ctx context.Context, _ selenium.WebDriver, _ *Args, version *po.Version, _ *sync.WaitGroup) (output *Args, exit Result, err error) {
	if version.Previous == nil {
		logs.InfoM(ctx, "check version skipped: no previous version")
		exit = Skipped
		return
	}
	previousVersion, err := version_util.Parse(a.VersionFormat, *version.Previous)
	if err != nil {
		logs.Error(ctx, "version parsing failed", err, "format", a.VersionFormat, "previous", *version.Previous)
		return
	}
	currentVersion, err := version_util.Parse(a.VersionFormat, version.Version)
	if err != nil {
		logs.Error(ctx, "version parsing failed", err, "format", a.VersionFormat, "current", version.Version)
		return
	}
	if previousVersion.EQ(currentVersion) {
		logs.InfoM(ctx, "check version interrupted: version not greater",
			"current", currentVersion, "previous", previousVersion)
		exit = StopFlow
	}
	logs.InfoM(ctx, "version update", "previousVersion", previousVersion, "currentVersion", currentVersion)
	return
}

func (a *CheckVersion) ToDTO() *DTO {
	return &DTO{
		Values: map[string]string{"format": a.VersionFormat},
	}
}
