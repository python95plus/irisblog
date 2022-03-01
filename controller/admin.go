package controller

import (
	"irisblog/config"
	"irisblog/middleware"
	"irisblog/provider"
	"irisblog/request"

	"github.com/kataras/iris/v12"
)

func AdminLogin(ctx iris.Context) {
	ctx.View("admin/login.html")
}

func AdminLoginForm(ctx iris.Context) {
	var user request.Admin
	if err := ctx.ReadForm(&user); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	admin, err := provider.GetAdminByUserName(user.UserName)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	if !admin.CheckPassword(user.Password) {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  "登录失败",
		})
		return
	}
	session := middleware.Sess.Start(ctx)
	session.Set("hasLogin", true)

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "登录成功",
		"data": 1,
	})
}

func AdminLogout(ctx iris.Context) {
	session := middleware.Sess.Start(ctx)
	session.Delete("hasLogin")
	ctx.Redirect("/")
}
