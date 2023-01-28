package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"software_updater/core/logs"
	"software_updater/core/util"
	"software_updater/ui/common"
)

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
