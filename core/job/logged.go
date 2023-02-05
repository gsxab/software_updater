package job

import (
	"github.com/tebeka/selenium"
	"golang.org/x/net/context"
	"software_updater/core/action"
	"software_updater/core/db/po"
	"software_updater/core/util/error_util"
	"sync"
)

type LoggedJob struct {
	DefaultJob
	info *DebugInfo
}

func (j *LoggedJob) RunAction(ctx context.Context, driver selenium.WebDriver, args *action.Args, v *po.Version, errs error_util.Collector, wg *sync.WaitGroup) (*action.Args, bool, bool, error) {
	output, stop, cancel, err := j.DefaultJob.RunAction(ctx, driver, args, v, errs, wg)
	j.info = &DebugInfo{
		Err:    err,
		Input:  args,
		Output: output,
	}
	return output, stop, cancel, err
}

func (j *LoggedJob) ToDTO() *DTO {
	return &DTO{
		State:     j.state.Int(),
		StateDesc: j.stateDesc,
		DebugInfo: j.info,
	}
}
