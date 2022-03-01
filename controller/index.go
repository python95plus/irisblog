package controller

import (
	"fmt"
	"irisblog/provider"

	"github.com/kataras/iris/v12"
)

func IndexPage(ctx iris.Context) {
	articles, _ := provider.GetArticleList(0, "id desc", 0, 10)
	categories, _ := provider.GetCategories(0)
	for num, article := range articles {
		for _, category := range categories {
			if article.CategoryId == category.Id {
				articles[num].Category = category
			}
		}
	}
	ctx.ViewData("articles", articles)
	ctx.View("index/index.html")
}

func ParamFunc(ctx iris.Context) {
	params := ctx.Params().GetEntry("params").Value().(string)
	fmt.Println("end...." + params)
	// for k, v := range params {
	// 	fmt.Println(k, v)
	// }
	ctx.WriteString(params)
}
