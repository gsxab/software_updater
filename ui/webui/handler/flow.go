/*
 * SPDX-License-Identifier: GPL-3.0-or-later
 *
 * Copyright (c) 2023. gsxab.
 *
 * This file is part of Software Update Watcher, a.k.a. Zhixin Robot.
 *
 * Software Update Watcher is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
 *
 * Software Update Watcher is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along with Software Update Watcher. If not, see <https://www.gnu.org/licenses/>.
 */

package handler

import (
	"html/template"
	"net/http"
	"software_updater/core/flow"
	"software_updater/core/util"
	"software_updater/ui/common"
	"strings"

	ginI18n "github.com/fishjar/gin-i18n"
	"github.com/gin-gonic/gin"
	"github.com/gsxab/logs"
)

type GetFlowRequest struct {
	Name   string `uri:"name" json:"name" form:"name" query:"name"`
	Reload bool   `json:"reload,omitempty" form:"reload" query:"reload"`
}

func GetFlow(ctx *gin.Context) {
	req := &GetFlowRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}
	if err := ctx.ShouldBindUri(req); err != nil {
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

	addFlowI18NInfo(ctx, data.BranchDTO)

	ctx.JSON(http.StatusOK, data)
}

const ActionNameI18NPrefix = "action_name."
const ActionHelpI18NPrefix = "action_help."

func addFlowI18NInfo(ctx *gin.Context, data *flow.BranchDTO) {
	for _, step := range data.Steps {
		actionKey := ActionNameI18NPrefix + step.ActionDTO.Name
		step.ActionDTO.I18NName = ginI18n.Msg(ctx, actionKey)
		actionHelpKey := ActionHelpI18NPrefix + step.ActionDTO.Name
		actionHelpValue := ginI18n.Msg(ctx, actionHelpKey)
		if actionHelpValue != actionHelpKey {
			tmpl, err := template.New(actionHelpKey).Parse(actionHelpValue)
			if err != nil {
				logs.Error(ctx, "template invalid", err, "key", actionHelpKey, "template", actionHelpValue)
				continue
			}
			sb := &strings.Builder{}
			err = tmpl.Execute(sb, step.ActionDTO.Values)
			if err != nil {
				logs.Error(ctx, "template invalid", err, "key", actionHelpKey, "template", actionHelpValue, "data", step.ActionDTO.Values)
				continue
			}
			step.ActionDTO.I18NHelp = sb.String()
		}
	}

	for _, branchDTO := range data.Next {
		addFlowI18NInfo(ctx, branchDTO)
	}
}
