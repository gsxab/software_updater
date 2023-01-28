package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"software_updater/core/logs"
	"software_updater/core/util"
	"software_updater/ui/common"
)

type GetFlowRequest struct {
	Name string
}

func GetFlow(ctx *gin.Context) {
	req := &GetFlowRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}

	data, err := common.GetFlowByName(ctx, req.Name, false)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, data)
}
