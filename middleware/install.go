package middleware

import (
	"irisblog/config"
	"strings"

	"github.com/kataras/iris/v12"
)

func CheckInstall(ctx iris.Context) {
	urlPath := ctx.RequestPath(false)
	if config.DB == nil && urlPath != "/install" && !strings.HasPrefix(urlPath, "/static") {
		ctx.Redirect("install")
		return
	} else {
		ctx.Next()
	}
}
