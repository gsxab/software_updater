package common

import (
	"context"
	"software_updater/core/db"
	"software_updater/core/db/po"
	"software_updater/core/logs"
	"software_updater/core/util"
	"software_updater/ui/dto"
)

func GetList(ctx context.Context, dateFormat string) ([]*dto.ListItemDTO, error) {
	//hpDAO := dao.Homepage
	//cvDAO := dao.CurrentVersion
	//vDAO := dao.Version

	hps := make([]*po.Homepage, 0)
	result := db.DB().Debug().WithContext(ctx).Model(&hps).
		Preload("Current").
		Preload("Current.Version").
		Find(&hps)
	if result.Error != nil {
		logs.Error(ctx, "list query failed", result.Error)
		return nil, result.Error
	}

	data := make([]*dto.ListItemDTO, 0, len(hps))
	for _, hp := range hps {
		datum := &dto.ListItemDTO{
			Name:    hp.Name,
			PageURL: hp.HomepageURL,
		}
		if hp.Current != nil {
			datum.ScheduledDate = util.FormatTime(hp.Current.ScheduledAt, dateFormat)
			if hp.Current.Version != nil {
				datum.Version = &hp.Current.Version.Version
				datum.UpdateDate = util.FormatTime(hp.Current.Version.LocalTime, dateFormat)
			}
		}
		data = append(data, datum)
	}

	return data, nil
}
