//go:build !core

package web

import (
	"embed"
	"io/fs"

	"github.com/gdy666/lucky/thirdlib/gdylib/ginutils"
	"github.com/gin-gonic/gin"
)

//go:embed adminviews/dist
var staticFs embed.FS
var stafs fs.FS

func initUI() {
	stafs, _ = fs.Sub(staticFs, "adminviews/dist")
}

func registerUIRoute(r *gin.Engine) {
	if stafs != nil {
		r.Use(ginutils.HandlerStaticFiles(stafs))
	}
}
