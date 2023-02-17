package base

import (
	"context"
	"software_updater/core/action"
	"software_updater/core/logs"
	"software_updater/core/util/version_util"
)

type VersionFilter struct {
	VersionFormat string `json:"format"`
	IndexReader
}

func (a *VersionFilter) Icon() string {
	return "mdi:mdi-alpha-v-circle"
}

func (a *VersionFilter) Filter(ctx context.Context, input *action.Args,
	callback func(v *version_util.Version) (bool, action.Result),
) (output *action.Args, exit action.Result, err error) {
	output = input

	versionStr, err := a.ReadDirectly(ctx, input)
	if err != nil {
		return
	}

	currentVersion, err := version_util.Parse(a.VersionFormat, versionStr)
	if err != nil {
		logs.Error(ctx, "version parsing failed", err, "format", a.VersionFormat, "versionStr", versionStr)
		return
	}
	res, exit := callback(currentVersion)
	if exit == action.Skipped {
		return
	}
	if !res {
		logs.InfoM(ctx, "version checker stopping task", "current", currentVersion)
		exit = action.StopFlow
	}
	return
}
