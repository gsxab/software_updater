package webui

import (
	ginI18n "github.com/fishjar/gin-i18n"
	"github.com/gin-gonic/gin"
	"io/fs"
	"mime"
	"net/http"
	"software_updater/core/config"
	config2 "software_updater/ui/webui/config"
	h "software_updater/ui/webui/handler"
	"software_updater/ui/webui/midware"
)

func RegisterRouters(r *gin.Engine) {
	ginI18n.LocalizerInit("zh-Hans", "zh-Hans,en-US", "./localize")
	r.Use(ginI18n.GinLocalizer())

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
	g.POST("/flow/start", h.StartFlow)
	g.POST("/flow/start_all", h.StartAllFlows)
	//g.POST("/flow/cancel", h.CancelFlow)
}
