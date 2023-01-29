package common

import (
	"context"
	"software_updater/core/db/dao"
	"software_updater/core/db/po"
	"software_updater/core/engine"
	"software_updater/core/logs"
)

func StartFlowByName(ctx context.Context, name string) (engine.TaskID, error) {
	hpDAO := dao.Homepage
	hp, err := hpDAO.WithContext(ctx).Where(hpDAO.Name.Eq(name)).Take()
	if err != nil {
		logs.Error(ctx, "homepage query failed", err, "name", name)
		return 0, err
	}
	data, err := StartFlow(ctx, hp)
	return data, err
}

func StartFlow(ctx context.Context, hp *po.Homepage) (engine.TaskID, error) {
	id, err := engine.Instance().Run(ctx, hp)
	if err != nil {
		return 0, err
	}
	return id, nil
}
