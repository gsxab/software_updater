package webui

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"software_updater/core/logs"
	"software_updater/core/util"
	"software_updater/ui/common"
)

func checkRPCSecret(ctx *gin.Context) {
	if webUIConfig.Secret != nil && ctx.Request.Header.Get("X-Soft-Token") != *webUIConfig.Secret {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
}

type GetListRequest struct {
}

func GetList(ctx *gin.Context) {
	req := &GetListRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}

	data, err := common.GetList(ctx, webUIConfig.DateFormat)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, data)
}

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

	data, err := common.GetVersionDetail(ctx, req.Name, nil, req.Version, webUIConfig.DateFormat)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, data)
}

type GetActionTreeRequest struct {
}

func GetActionTree(ctx *gin.Context) {
	req := &GetActionTreeRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}

	data, err := common.GetActionTree(ctx)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, data)
}
