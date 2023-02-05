package std

import (
	"context"
	"github.com/tebeka/selenium"
	"software_updater/core/action"
	"software_updater/core/action/base"
	"software_updater/core/db/po"
	"software_updater/core/util/version_util"
	"sync"
)

type CheckLaterVersion struct {
	base.DefaultFactory[CheckLaterVersion, *CheckLaterVersion]
	base.VersionComparer
}

func (a *CheckLaterVersion) Path() action.Path {
	return action.Path{"basic", "value_check", "version_gt"}
}

func (a *CheckLaterVersion) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, version *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.Compare(ctx, input, version.Previous, func(prevV *version_util.Version, newV *version_util.Version) bool {
		return prevV.LT(newV)
	})
}

func (a *CheckLaterVersion) ToDTO() *action.DTO {
	return &action.DTO{
		Values: map[string]string{"format": a.VersionFormat},
	}
}
