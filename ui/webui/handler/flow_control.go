package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"software_updater/core/logs"
	"software_updater/core/util"
	"software_updater/ui/common"
)

type StartFlowRequest struct {
	Name string
}

type StartFlowData struct {
	ID int64
}

func StartFlow(ctx *gin.Context) {
	req := &StartFlowRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}

	id, err := common.StartFlowByName(ctx, req.Name)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, &StartFlowData{ID: id})
}

//type CancelFlowRequest struct {
//	Name string
//}
//
//func CancelFlow(ctx *gin.Context) {
//	req := &CancelFlowRequest{}
//	if err := ctx.ShouldBind(req); err != nil {
//		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
//		ctx.Status(http.StatusBadRequest)
//		return
//	}
//
//	data, err := common.CancelFlowByName(ctx, req.Name, false)
//	if err != nil {
//		ctx.Status(http.StatusInternalServerError)
//		return
//	}
//
//	ctx.JSON(http.StatusOK, data)
//}
