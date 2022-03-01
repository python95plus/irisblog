package controller

import "github.com/kataras/iris/v12"

func InternalServerError(ctx iris.Context) {
	errMessage := ctx.Values().GetString("error")
	if errMessage == "" {
		errMessage = "(Unexpected) internal server error"
	}

	ctx.ViewData("errMessage", errMessage)
	ctx.View("errors/500.html")
}
