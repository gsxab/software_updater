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

package webui

import (
	"io/fs"
	"mime"
	"net/http"
	"software_updater/core/config"
	config2 "software_updater/ui/webui/config"
	h "software_updater/ui/webui/handler"
	"software_updater/ui/webui/midware"

	ginI18n "github.com/fishjar/gin-i18n"
	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine) {
	r.Use(ginI18n.Localizer(&ginI18n.Options{
		DefaultLang:  "zh-Hans",
		SupportLangs: "zh-Hans,en-US",
		FilePath:     "./localize",
	}))

	// screenshots
	r.StaticFS("/static/screenshot", http.Dir(config.Current().Files.ScreenshotDir))
	// generated html (compat)
	//r.StaticFS("/static/html", http.Dir(config.Current().Files.HTMLDir))
	// vue ui
	jsFS, _ := fs.Sub(DistFiles, "dist/static/js")
	cssFS, _ := fs.Sub(DistFiles, "dist/static/css")
	fontFS, _ := fs.Sub(DistFiles, "dist/static/fonts")
	//r.StaticFS("/static/", http.FS(staticFS))
	r.StaticFS("/static/js/", http.FS(jsFS))
	r.StaticFS("/static/css/", http.FS(cssFS))
	r.StaticFS("/static/fonts/", http.FS(fontFS))
	r.GET("/", func(ctx *gin.Context) {
		//	ctx.Redirect(http.StatusMovedPermanently, "/index.html")
		//})
		//r.GET("/index.html", func(ctx *gin.Context) {
		data, err := DistFiles.ReadFile("dist/index.html")
		if err != nil {
		}
		ctx.Data(http.StatusOK, mime.TypeByExtension(".html"), data)
	})
	r.GET("/favicon.ico", func(ctx *gin.Context) {
		data, err := DistFiles.ReadFile("dist/favicon.ico")
		if err != nil {
		}
		ctx.Data(http.StatusOK, mime.TypeByExtension(".ico"), data)
	})

	// rpc
	g := r.Group("/jsonrpc/v1", midware.CheckRPCSecret(config2.WebUIConfig.Secret))

	// list
	g.GET("/list", h.GetList)

	// version
	g.GET("/version/:name/:version", h.GetVersion)

	// action
	g.GET("/actions", h.GetActionTree)

	// flow
	g.GET("/flow/:name", h.GetFlow)
	//g.GET("/flow/:name/edit", h.GetFlowEditForm)
	//g.POST("/flow/:name", EditFlow)
	//g.GET("/flow/:name/create", h.GetFlowCreateForm)
	//g.PUT("/flow/:name", CreateFlow)

	// start task from flow
	g.POST("/flow/:name/start", h.StartFlow)
	g.POST("/flow/all/start", h.StartAllFlows)

	// task
	g.GET("/tasks", h.GetTaskIDMap)
	//g.POST("/task/:id/cancel", h.CancelTask)
	//g.POST("/task/all/cancel", h.CancelAllTasks)
	g.GET("/task/:id", h.GetTask)
}
