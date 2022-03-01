package controller

import (
	"log"

	"irisblog/config"
	"irisblog/provider"
	"irisblog/request"

	"github.com/kataras/iris/v12"
)

func Installer(ctx iris.Context) {
	if config.DB != nil {
		ctx.Redirect("/")
		return
	}
	// dataNow := time.Now().Unix()
	// ctx.ViewData("datanow", dataNow)
	ctx.ViewData("SiteName", "博客信息配置")
	ctx.View("install/index.html")
}

func InstallForm(ctx iris.Context) {
	// blog host port database user password admin_user admin_password
	var formData request.Install
	if errForm := ctx.ReadForm(&formData); errForm != nil {
		log.Println(errForm.Error())
	}

	config.JsonData.Server.SiteName = formData.Blog
	config.JsonData.DB.Database = formData.Database
	config.JsonData.DB.Host = formData.Host
	config.JsonData.DB.Port = formData.Port
	config.JsonData.DB.User = formData.User
	config.JsonData.DB.Password = "" // password

	errDB := config.InitDB(&config.JsonData.DB)
	if errDB != nil {
		log.Println(errDB.Error())
	}

	if errAdmin := provider.InitAdmin(formData.AdminUser, formData.AdminPassword); errAdmin != nil {
		log.Println(errAdmin.Error())
	}
	err := config.WriteConfig()
	if err != nil {
		log.Println(err.Error())
	}
	ctx.Redirect("/")
}
