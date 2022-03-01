package route

import (
	"fmt"
	"irisblog/controller"
	"regexp"
	"strings"

	"github.com/kataras/iris/v12"
)

func Register(app *iris.Application) {
	registerMacros(app)
	app.Get("/", controller.IndexPage).Name = "index"
	app.Get("/install", controller.Installer).Name = "install"
	app.Post("/install", controller.InstallForm).Name = "installform"
	app.Get("/{params:rewrite}", controller.ParamFunc)
	admin := app.Party("/admin")
	{
		admin.Get("/login", controller.AdminLogin).Name = "login"
		admin.Post("/login", controller.AdminLoginForm)
		admin.Get("/logout", controller.AdminLogout).Name = "logout"
	}

	article := app.Party("article")
	{
		article.Get("/{id:rewrite}", controller.ArticleDetail).Name = "article"
		article.Get("/publish", controller.ArticlePublish).Name = "article_publish"
		article.Post("/publish", controller.ArticlePublishForm)
	}

	category := app.Party("category")
	{
		category.Get("/{id:uint}", controller.CategoryList).Name = "category"
	}

	attachment := app.Party("/attachment")
	{
		attachment.Post("/upload", controller.AttachmentUpload)
	}
}

var reg = regexp.MustCompile(`\d+$`)

func registerMacros(app *iris.Application) {
	app.Macros().Register("rewrite", "", false, true, func(paramValue string) (interface{}, bool) {
		switch {
		case strings.HasSuffix(paramValue, ".html"):
			paramValue = strings.Replace(paramValue, ".html", "", 1)
		case strings.Contains(paramValue, "."):
			return nil, false
		case reg.MatchString(paramValue):
			fmt.Println("bbbbbbbbbbbbb")
		}
		return paramValue, true
	})
}
