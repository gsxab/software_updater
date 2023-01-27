package common

import (
	"context"
	"software_updater/core/engine"
	"software_updater/ui/dto"
)

func GetActionTree(ctx context.Context) (*dto.ActionHierarchyDTO, error) {
	data, err := engine.Instance().ActionHierarchy(ctx)
	return data, err
}
