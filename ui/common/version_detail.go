package common

import (
	"context"
	"software_updater/core/db/dao"
	"software_updater/core/logs"
	"software_updater/core/util"
	"software_updater/core/util/optional"
	"software_updater/ui/dto"
)

func GetVersionDetail(ctx context.Context, name string, optionalPage *string, v string, dateFormat string) (*dto.VersionDTO, error) {
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
		PrevVersion: version.Previous,
		NextVersion: nil,
		RemoteDate:  util.FormatTime(version.RemoteDate, dateFormat),
		UpdateDate:  *util.FormatTime(version.LocalTime, dateFormat),
		Link:        version.Link,
		Digest:      version.Digest,
		Picture:     version.Picture,
	}

	nextVersionExist, err := vDAO.WithContext(ctx).Where(vDAO.Name.Eq(name), vDAO.Previous.Eq(v)).Count()
	if err != nil {
		logs.Error(ctx, "version query failed", err, "name", name, "v", v)
		return nil, err
	}

	if nextVersionExist > 0 {
		nextVersion, err := vDAO.WithContext(ctx).Where(vDAO.Name.Eq(name), vDAO.Previous.Eq(v)).Take()
		if err != nil {
			logs.Error(ctx, "version query failed", err, "name", name, "v", v)
			return nil, err
		}

		data.NextVersion = &nextVersion.Version
	}

	return data, nil
}
