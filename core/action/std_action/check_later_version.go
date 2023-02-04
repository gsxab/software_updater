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

type CheckLaterVersion struct {
	action.Default
	action.DefaultFactory[CheckLaterVersion, *CheckLaterVersion]
	VersionFormat string `json:"format"`
	action.IndexReader
}

func (a *CheckLaterVersion) Path() action.Path {
	return action.Path{"basic", "value_check", "version_gt"}
}

func (a *CheckLaterVersion) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, version *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	versionStr, err := a.ReadDirectly(ctx, input)
	if err != nil {
		return
	}

	if version.Previous == nil {
		logs.InfoM(ctx, "check version skipping: no previous version")
		exit = action.Skipped
		return
	}
	previousVersion, err := version_util.Parse(a.VersionFormat, version.Version)
	if err != nil {
		logs.Error(ctx, "version parsing failed", err, "format", a.VersionFormat, "previous", version.Version)
		return
	}

	currentVersion, err := version_util.Parse(a.VersionFormat, versionStr)
	if err != nil {
		logs.Error(ctx, "version parsing failed", err, "format", a.VersionFormat, "previous", version.Version)
		return
	}
	if currentVersion.LE(currentVersion) {
		logs.InfoM(ctx, "check version exiting: version not greater",
			"current", currentVersion, "previous", previousVersion)
		exit = action.StopFlow
	}
	logs.InfoM(ctx, "version update", "previousVersion", previousVersion, "currentVersion", currentVersion)
	return
}

func (a *CheckLaterVersion) ToDTO() *action.DTO {
	return &action.DTO{
		Values: map[string]string{"format": a.VersionFormat},
	}
}
