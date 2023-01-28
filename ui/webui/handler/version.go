package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"software_updater/core/logs"
	"software_updater/core/util"
	"software_updater/ui/common"
	"software_updater/ui/webui/config"
)

type GetVersionRequest struct {
	Name    string `json:"name" form:"name" query:"name"`
	Version string `json:"version" form:"version" query:"version"`
}

func GetVersion(ctx *gin.Context) {
	req := &GetVersionRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}

	data, err := common.GetVersionDetail(ctx, req.Name, nil, req.Version, config.WebUIConfig.DateFormat)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, data)
}
