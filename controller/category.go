package controller

import (
	"irisblog/provider"

	"github.com/kataras/iris/v12"
)

func CategoryList(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	articles, _ := provider.GetArticleList(id, "id desc", 0, 10)
	ctx.ViewData("articles", articles)
	ctx.View("index/index.html")
}
