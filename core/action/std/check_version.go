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

type CheckVersion struct {
	base.DefaultFactory[CheckVersion, *CheckVersion]
	base.VersionComparer
}

func (a *CheckVersion) Path() action.Path {
	return action.Path{"basic", "value_check", "version_neq"}
}

func (a *CheckVersion) Do(ctx context.Context, _ selenium.WebDriver, input *action.Args, version *po.Version, _ *sync.WaitGroup) (output *action.Args, exit action.Result, err error) {
	return a.Compare(ctx, input, version.Previous, func(prevV *version_util.Version, newV *version_util.Version) bool {
		return !newV.EQ(prevV)
	})
}

func (a *CheckVersion) ToDTO() *action.DTO {
	return &action.DTO{
		Values: map[string]string{"format": a.VersionFormat},
	}
}
