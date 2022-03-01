package main

import (
	"fmt"
	"irisblog/config"
	"irisblog/middleware"
	"irisblog/route"
	"irisblog/tags"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

type Bootstrap struct {
	Application *iris.Application
	Port        int
	LoggerLevel string
}

func NewApp(port int, loggerLevel string) *Bootstrap {
	var bootstrap Bootstrap
	bootstrap.Application = iris.New()
	bootstrap.Port = port
	bootstrap.LoggerLevel = loggerLevel
	return &bootstrap
}

func (bootstrap *Bootstrap) LoadGlobalMiddleware() {
	bootstrap.Application.Use(logger.New())
	bootstrap.Application.Use(recover.New())
	bootstrap.Application.Use(middleware.Auth)
	bootstrap.Application.Use(middleware.CheckInstall)
}

func (bootstrap *Bootstrap) Serve() {
	if config.DB == nil {
		fmt.Println("DB no value")
	}
	// bootstrap.Application.Logger().SetLevel(bootstrap.LoggerLevel)
	bootstrap.LoadGlobalMiddleware()

	bootstrap.Application.HandleDir("/static", "./static")
	pugEngine := iris.Django("./templates", ".html")

	if config.SeverConfig.Env == "development" {
		pugEngine.Reload(true)
	}

	pugEngine.AddFunc("timestampToDate", TimestampToDate)
	pugEngine.AddFunc("homeTitle", HomeTitle)
	pugEngine.RegisterTag("categoryList", tags.TagCategoryListParser)
	pugEngine.RegisterTag("articleList", tags.TagArticleListParser)

	bootstrap.Application.RegisterView(pugEngine)
	route.Register(bootstrap.Application) //定义router路由信息
	bootstrap.Application.Run(
		iris.Addr(fmt.Sprintf(":%d", bootstrap.Port)),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithoutBodyConsumptionOnUnmarshal,
	)

}

func TimestampToDate(in int64, layout string) string {
	t := time.Unix(in, 0)
	return t.Format(layout)
}

func HomeTitle() string {
	return config.SeverConfig.SiteName
}

func main() {
	app := NewApp(config.SeverConfig.Port, config.SeverConfig.Env)
	app.Serve()
}
