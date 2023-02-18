package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"software_updater/core/logs"
	"software_updater/core/util"
	"software_updater/ui/common"
)

type StartFlowRequest struct {
	Name string `json:"name" form:"name" query:"name"`
}

type StartAllFlowsRequest struct {
	Force bool `json:"force" form:"force" query:"force"`
}

type StartFlowData struct {
	ID int64 `json:"id"`
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

func StartAllFlows(ctx *gin.Context) {
	req := &StartAllFlowsRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}

	err := common.StartAllFlows(ctx)
	if err != nil {
		logs.Error(ctx, "start all flows failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
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
