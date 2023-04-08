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
	"github.com/gin-gonic/gin"
	"net/http"
	"software_updater/core/logs"
	"software_updater/core/util"
	"software_updater/ui/common"
)

type StartFlowRequest struct {
	Name  string `uri:"name" form:"name" query:"name"`
	Force bool   `json:"force,omitempty" form:"force" query:"force"`
}

type StartFlowData struct {
	ID int64 `json:"id"`
}

type StartAllFlowData struct {
	IDs []int64 `json:"ids"`
}

func StartFlow(ctx *gin.Context) {
	req := &StartFlowRequest{}
	if err := ctx.ShouldBindUri(req); err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}

	idMap, err := common.StartFlow(ctx, req.Name, req.Force)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if len(idMap) != 1 {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	for _, id := range idMap {
		ctx.JSON(http.StatusCreated, &StartFlowData{ID: id})
	}
}

func StartAllFlows(ctx *gin.Context) {
	req := &StartFlowRequest{}
	if err := ctx.ShouldBindUri(req); err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}
	req.Name = "all"

	idMap, err := common.StartFlow(ctx, req.Name, req.Force)
	if err != nil {
		if req.Name == "all" {
			ids := make([]int64, 0, len(idMap))
			for _, id := range idMap {
				ids = append(ids, id)
			}
			ctx.JSON(http.StatusInternalServerError, &StartAllFlowData{IDs: ids})
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	if req.Name == "all" {
		ids := make([]int64, 0, len(idMap))
		for _, id := range idMap {
			ids = append(ids, id)
		}
		ctx.JSON(http.StatusCreated, &StartAllFlowData{IDs: ids})
		return
	}

	if len(idMap) != 1 {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	for _, id := range idMap {
		ctx.JSON(http.StatusCreated, &StartFlowData{ID: id})
	}
}
