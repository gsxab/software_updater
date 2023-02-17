package webui

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"software_updater/core/logs"
	"software_updater/ui/webui/config"
	"syscall"
	"time"
)

func InitAndRun(ctx context.Context, configExtraUI string) error {
	config.WebUIConfig = config.DefaultConfig()
	if configExtraUI != "" {
		err := json.Unmarshal([]byte(configExtraUI), config.WebUIConfig)
		if err != nil {
			return err
		}
	}

	r := gin.Default()

	RegisterRouters(r)

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
