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
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"software_updater/ui/webui/config"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsxab/logs"
)

func InitAndRun(ctx context.Context, configExtraUI string) error {
	logs.InfoM(ctx, `
Software Update Watcher, a.k.a. Zhixin Robot  Copyright (C) 2023  gsxab
This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it under certain conditions.
`)

	config.WebUIConfig = config.DefaultConfig()
	if configExtraUI != "" {
		err := json.Unmarshal([]byte(configExtraUI), config.WebUIConfig)
		if err != nil {
			return err
		}
	}

	r := gin.Default()

	RegisterRouters(ctx, r)

	srv := &http.Server{
		Addr:         config.WebUIConfig.Addr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 8)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	logs.WarnM(ctx, "shutdown received")
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.Error(ctx, "shutdown failed", err)
		syscall.Exit(1)
	}
	<-ctx.Done()
	logs.InfoM(ctx, "exiting")

	return nil
}
