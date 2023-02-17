package handler

import (
	ginI18n "github.com/fishjar/gin-i18n"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"software_updater/core/job"
	"software_updater/core/logs"
	"software_updater/core/util"
	"software_updater/ui/common"
	"strings"
)

type GetFlowRequest struct {
	Name   string `json:"name" form:"name" query:"name"`
	Reload bool   `json:"reload,omitempty" form:"reload" query:"reload"`
}

func GetFlow(ctx *gin.Context) {
	req := &GetFlowRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}

	data, err := common.GetFlowByName(ctx, req.Name, req.Reload)
	if err != nil {
		logs.ErrorE(ctx, err, "name", req.Name)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	localizer := ctx.MustGet("Localizer").(*ginI18n.UserLocalize)
	addFlowI18NInfo(ctx, localizer, data.BranchDTO)

	ctx.JSON(http.StatusOK, data)
}

const ActionNameI18NPrefix = "action_name."
const ActionHelpI18NPrefix = "action_help."

func addFlowI18NInfo(ctx *gin.Context, l *ginI18n.UserLocalize, data *job.BranchDTO) {
	for _, job := range data.Jobs {
		actionKey := ActionNameI18NPrefix + job.ActionDTO.Name
		job.ActionDTO.I18NName = l.GetMsg(actionKey)
		actionHelpKey := ActionHelpI18NPrefix + job.ActionDTO.Name
		actionHelpValue := l.GetMsg(actionHelpKey)
		if actionHelpValue != actionHelpKey {
			tmpl, err := template.New(actionHelpKey).Parse(actionHelpValue)
			if err != nil {
				logs.Error(ctx, "template invalid", err, "key", actionHelpKey, "template", actionHelpValue)
				continue
			}
			sb := &strings.Builder{}
			err = tmpl.Execute(sb, job.ActionDTO.Values)
			if err != nil {
				logs.Error(ctx, "template invalid", err, "key", actionHelpKey, "template", actionHelpValue, "data", job.ActionDTO.Values)
				continue
			}
			job.ActionDTO.I18NHelp = sb.String()
		}
	}

	for _, branchDTO := range data.Next {
		addFlowI18NInfo(ctx, l, branchDTO)
	}
}
