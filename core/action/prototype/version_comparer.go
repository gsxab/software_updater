package prototype

import (
	"context"
	"software_updater/core/action"
	"software_updater/core/db/dao"
	"software_updater/core/logs"
	"software_updater/core/util/version_util"
)

type VersionComparer struct {
	VersionFilter
}

func (a *VersionComparer) Compare(ctx context.Context, input *action.Args, prevID *uint,
	callback func(prevV *version_util.Version, newV *version_util.Version) bool,
) (output *action.Args, exit action.Result, err error) {
	return a.Filter(ctx, input, func(newVersion *version_util.Version) (res bool, exit action.Result) {
		if prevID == nil {
			logs.InfoM(ctx, "check version skipping: no previous version")
			exit = action.Skipped
		}

		vDAO := dao.Version
		prev, err := vDAO.WithContext(ctx).Where(vDAO.ID.Eq(*prevID)).Take()
		if err != nil {
			logs.Error(ctx, "version query failed", err, "id", *prevID)
			return
		}

		previousVersion, err := version_util.Parse(a.VersionFormat, prev.Version)
		if err != nil {
			logs.Error(ctx, "version parsing failed", err, "format", a.VersionFormat, "previous", prev.Version)
			return
		}
		res = callback(previousVersion, newVersion)
		return
	})
}
