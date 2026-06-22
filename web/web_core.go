//go:build core

package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func initUI() {
	// core 模式不加载本地前端
}

func registerUIRoute(r *gin.Engine) {
	// core 模式将根路径重定向到公网面板
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "http://lucky666.cn/admin/")
	})
}
