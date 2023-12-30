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
	"net/http"
	"software_updater/core/util"
	"software_updater/ui/common"
	"software_updater/ui/dto"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gsxab/go-logs"
)

type GetTaskStateRequest struct {
	TaskID int64 `uri:"id"`
}

type GetTaskStateData struct {
	Exist bool `json:"exist"`
	State int  `json:"state"`
}

func GetTaskState(ctx *gin.Context) {
	req := &GetTaskStateRequest{}
	if err := ctx.ShouldBindUri(req); err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}

	exist, state, err := common.GetTaskStateByID(ctx, req.TaskID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, &GetTaskStateData{Exist: exist, State: state})
}

type GetTaskIDMapRequest struct {
}

func GetTaskIDMap(ctx *gin.Context) {
	req := &GetTaskIDMapRequest{}
	if err := ctx.ShouldBindUri(req); err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}

	idMap, err := common.GetTaskIDMap(ctx)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, idMap)
}

type GetTaskRequest struct {
	TaskID string `uri:"id"`
}

func GetTask(ctx *gin.Context) {
	req := &GetTaskRequest{}
	if err := ctx.ShouldBindUri(req); err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}

	taskID, err := strconv.ParseInt(req.TaskID, 10, 64)
	if err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", req, util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}
	exists, data, err := common.GetTaskMeta(ctx, taskID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if !exists {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

type GetTaskListRequest struct {
	TaskIDs string `uri:"ids"`
}

func GetTaskList(ctx *gin.Context) {
	req := &GetTaskListRequest{}
	if err := ctx.ShouldBindUri(req); err != nil {
		logs.Warn(ctx, "request param resolving failed", err, "req", util.ToJSON(req))
		ctx.Status(http.StatusBadRequest)
		return
	}

	if req.TaskIDs != "" {
		splitIDs := strings.Split(req.TaskIDs, ",")
		data := make([]*dto.TaskMetaDTO, 0, len(splitIDs))
		for _, taskIDStr := range splitIDs {
			taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
			if err != nil {
				logs.Warn(ctx, "request param resolving failed", err, "req", req, util.ToJSON(req))
				ctx.Status(http.StatusBadRequest)
				return
			}
			exists, taskMeta, err := common.GetTaskMeta(ctx, taskID)
			if err != nil {
				ctx.Status(http.StatusInternalServerError)
				return
			}
			if !exists {
				ctx.Status(http.StatusNotFound)
				return
			}
			data = append(data, taskMeta)
		}
		ctx.JSON(http.StatusOK, data)
	}

	data, err := common.GetTaskMetaList(ctx)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, data)
}
