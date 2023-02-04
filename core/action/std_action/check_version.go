package std_action

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/util/version_util"
	"sync"
)

type CheckVersion struct {
	action.Default
	action.DefaultFactory[CheckVersion, *CheckVersion]
	VersionFormat string `json:"format"`
}

func (a *CheckVersion) Path() action.Path {
	return action.Path{"basic", "value_check", "version_neq"}
}

func (a *CheckVersion) Do(ctx context.Context, _ selenium.WebDriver, _ *action.Args, version *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	if version.Previous == nil {
		logs.InfoM(ctx, "check version skipped: no previous version")
		exit = action.Skipped
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
		exit = action.StopFlow
	}
	logs.InfoM(ctx, "version update", "previousVersion", previousVersion, "currentVersion", currentVersion)
	return
}

func (a *CheckVersion) ToDTO() *action.DTO {
	return &action.DTO{
		Values: map[string]string{"format": a.VersionFormat},
	}
}
