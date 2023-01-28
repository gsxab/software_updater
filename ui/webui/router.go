package webui

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"software_updater/core/config"
	config2 "software_updater/ui/webui/config"
	h "software_updater/ui/webui/handler"
	"software_updater/ui/webui/midware"
)

func RegisterRouters(r *gin.Engine) {
	// screenshots
	r.StaticFS("/static/screenshot", http.Dir(config.Current().Files.ScreenshotDir))
	// generated html
	r.StaticFS("/static/html", http.Dir(config.Current().Files.HTMLDir))

	// rpc
	g := r.Group("/jsonrpc/v1", midware.CheckRPCSecret(config2.WebUIConfig.Secret))
	// list
	g.GET("/list", h.GetList)
	// version
	g.GET("/version", h.GetVersion)
	// action
	g.GET("/action", h.GetActionTree)
	// flow
	g.GET("/flow", h.GetFlow)
	//g.GET("/flow/realtime", getRealTimeFlow)
	//g.GET("/flow/job", getJob)
	//g.PUT("/flow", putFlow)
	//g.DELETE("/flow/job", deleteFlow)
	// flow state
	//g.POST("/flow/start", startFlow)
	//g.POST("/flow/cancel", cancelFlow)
}
