package midware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckRPCSecret(secret *string) func(ctx *gin.Context) {
	if secret != nil {
		return func(ctx *gin.Context) {
			if ctx.Request.Header.Get("X-Soft-Token") != *secret {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}
		}
	} else {
		return func(ctx *gin.Context) {}
	}
}
