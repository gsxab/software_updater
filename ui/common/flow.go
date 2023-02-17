package common

import (
	"context"
	"software_updater/core/db/dao"
	"software_updater/core/db/po"
	"software_updater/core/engine"
	"software_updater/core/job"
	"software_updater/core/logs"
	"software_updater/ui/dto"
)

func GetFlowByName(ctx context.Context, name string, reload bool) (*dto.FlowDTO, error) {
	hpDAO := dao.Homepage
	hp, err := hpDAO.WithContext(ctx).Where(hpDAO.Name.Eq(name)).Take()
	if err != nil {
		logs.Error(ctx, "homepage query failed", err, "name", name)
		return nil, err
	}
	data, err := GetFlow(ctx, hp, reload)
	return data, err
}

func GetFlow(ctx context.Context, hp *po.Homepage, reload bool) (*dto.FlowDTO, error) {
	flow, err := engine.Instance().Load(ctx, hp, !reload)
	if err != nil {
		return nil, err
	}
	data := job.ToFlowDTO(flow)
	return data, nil
}
