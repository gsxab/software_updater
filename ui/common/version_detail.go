package common

import (
	"context"
	"software_updater/core/db/dao"
	"software_updater/core/logs"
	"software_updater/core/util/optional"
	"software_updater/ui/dto"
	"time"
)

func GetVersionDetail(ctx context.Context, name string, optionalPage *string, v string) (*dto.VersionDTO, error) {
	page, err := optional.OrLazy(optionalPage, func() (string, error) {
		hpDAO := dao.Homepage
		hp, err := hpDAO.WithContext(ctx).Where(hpDAO.Name.Eq(name)).Take()
		if err != nil {
			logs.Error(ctx, "version query failed", err, "name", name, "v", v)
			return "", err
		}
		return hp.HomepageURL, nil
	})
	if err != nil {
		return nil, err
	}

	vDAO := dao.Version

	version, err := vDAO.WithContext(ctx).Where(vDAO.Name.Eq(name), vDAO.Version.Eq(v)).Take()
	if err != nil {
		logs.Error(ctx, "version query failed", err, "name", name, "v", v)
		return nil, err
	}

	data := &dto.VersionDTO{
		Name:        name,
		HomepageURL: page,
		Version:     v,
		PrevVersion: nil,
		NextVersion: nil,
		RemoteDate:  dto.ToDateDTO(version.RemoteDate, time.UTC),
		UpdateDate:  dto.ToDateDTO(version.LocalTime, time.Local),
		Link:        version.Link,
		Digest:      version.Digest,
		Picture:     version.Picture,
	}

	if version.Previous != nil {
		previousVersion, err := vDAO.WithContext(ctx).Where(vDAO.ID.Eq(*version.Previous)).Take()
		if err != nil {
			logs.Error(ctx, "version query failed", err, "id", *version.Previous)
			return nil, err
		}

		data.PrevVersion = &previousVersion.Version
	}

	nextVersionSlice, err := vDAO.WithContext(ctx).Where(vDAO.Previous.Eq(version.ID)).Find()
	if err != nil {
		logs.Error(ctx, "version query failed", err, "name", name, "v", v)
		return nil, err
	}

	if len(nextVersionSlice) == 1 {
		data.NextVersion = &nextVersionSlice[0].Version
	} else if len(nextVersionSlice) > 1 {
		logs.ErrorM(ctx, "next version more than one", "id", version.ID)
	}

	return data, nil
}
