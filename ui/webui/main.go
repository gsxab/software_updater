package webui

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"software_updater/core/config"
	"software_updater/core/logs"
	"syscall"
	"time"
)

func InitAndRun(ctx context.Context, configExtraUI string) error {
	webUIConfig = DefaultConfig()
	if configExtraUI != "" {
		err := json.Unmarshal([]byte(configExtraUI), webUIConfig)
		if err != nil {
			return err
		}
	}

	r := gin.Default()

	// screenshots
	r.StaticFS("/static/screenshot", http.Dir(config.Current().Files.ScreenshotDir))
	// generated html
	r.StaticFS("/static/html", http.Dir(config.Current().Files.HTMLDir))

	// rpc
	g := r.Group("/jsonrpc/v1", checkRPCSecret)
	// list
	g.GET("/list", GetList)
	// version
	g.GET("/version", GetVersion)
	// action
	g.GET("/action", GetActionTree)
	// flow
	//g.GET("/flow", getFlow)
	//g.GET("/flow/realtime", getRealTimeFlow)
	//g.GET("/flow/job", getJob)
	//g.PUT("/flow", putFlow)
	//g.DELETE("/flow/job", deleteFlow)
	// flow state
	//g.POST("/flow/start", startFlow)
	//g.POST("/flow/cancel", cancelFlow)

	srv := &http.Server{
		Addr:         webUIConfig.Addr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	logs.WarnM(ctx, "shutdown received")
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logs.Error(ctx, "shutdown failed", err)
		syscall.Exit(1)
	}
	select {
	case <-ctx.Done():
	}
	logs.InfoM(ctx, "exiting")

	return nil
}
