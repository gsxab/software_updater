package common

import (
	"context"
	"software_updater/core/db/dao"
	"software_updater/core/logs"
	"software_updater/ui/dto"
	"time"
)

func GetList(ctx context.Context) ([]*dto.ListItemDTO, error) {
	hpDAO := dao.Homepage

	hps, err := hpDAO.WithContext(ctx).Preload(hpDAO.Current).Preload(hpDAO.Current.Version).Find()
	if err != nil {
		logs.Error(ctx, "list query failed", err)
		return nil, err
	}

	data := make([]*dto.ListItemDTO, 0, len(hps))
	for _, hp := range hps {
		datum := &dto.ListItemDTO{
			Name:    hp.Name,
			PageURL: hp.HomepageURL,
		}
		if hp.Current != nil {
			datum.ScheduledDate = dto.ToDateDTO(hp.Current.ScheduledAt, time.Local)
			if hp.Current.Version != nil {
				datum.Version = &hp.Current.Version.Version
				datum.UpdateDate = dto.ToDateDTO(hp.Current.Version.LocalTime, time.Local)
			}
		}
		data = append(data, datum)
	}

	return data, nil
}
