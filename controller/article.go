package controller

import (
	"fmt"
	"irisblog/config"
	"irisblog/model"
	"irisblog/provider"
	"irisblog/request"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kataras/iris/v12"
)

func ArticleDetail(ctx iris.Context) {
	id := ctx.Params().GetUintDefault("id", 0)
	article, err := provider.GetArticleById(id)
	if err != nil {
		return
	}
	hasLogin := ctx.Values().GetBoolDefault("hasLogin", false)
	ctx.ViewData("hasLogin", hasLogin)
	if !hasLogin {
		article.AddViews(config.DB)
	}
	prev, _ := provider.GetPrevArticleById(article.CategoryId, article.Id)
	next, _ := provider.GetNextArticleById(article.CategoryId, article.Id)
	ctx.ViewData("prev", prev)
	ctx.ViewData("next", next)
	ctx.ViewData("article", article)
	ctx.View("article/detail.html")
}

func ArticlePublish(ctx iris.Context) {
	if !ctx.Values().GetBoolDefault("hasLogin", false) {
		InternalServerError(ctx)
		return
	}
	id := uint(ctx.URLParamIntDefault("id", 0))
	if id > 0 {
		article, _ := provider.GetArticleById(id)
		ctx.ViewData("article", article)
	}
	categories, err := provider.GetCategories(0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	categoryList := make([]string, 0)

	for _, category := range categories {
		categoryList = append(categoryList, category.Title)
	}

	ctx.ViewData("categories", categoryList)

	ctx.View("/article/publish.html")
}

func ArticlePublishForm(ctx iris.Context) {
	var err error
	if !ctx.Values().GetBoolDefault("hasLogin", false) {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  "还未登录",
		})
		return
	}
	var req request.Article
	if err = ctx.ReadForm(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	var category *model.Category
	if req.CategoryName != "" {
		category, err = provider.GetCategoryByTitle(req.CategoryName)
		if err != nil {
			category = &model.Category{
				Title:  req.CategoryName,
				Status: 1,
			}
			err = category.Save(config.DB)
			if err != nil {
				ctx.JSON(iris.Map{
					"code": config.StatusFailed,
					"msg":  err.Error(),
				})
				return
			}
		}
	}

	var article *model.Article
	if req.Id > 0 {
		article, err = provider.GetArticleById(req.Id)
		if err != nil {
			ctx.JSON(iris.Map{
				"code": config.StatusFailed,
				"msg":  err.Error(),
			})
			return
		}
		if article.ArticleData == nil {
			article.ArticleData = &model.ArticleData{}
		}
		article.ArticleData.Content = req.Content
	} else {
		article = &model.Article{
			Title:       req.Title,
			Keywords:    req.Keywords,
			Description: req.Description,
			Status:      1,
			ArticleData: &model.ArticleData{
				Content: req.Content,
			},
		}
	}

	if req.Description == "" {
		htmlR := strings.NewReader(req.Content)
		doc, errdoc := goquery.NewDocumentFromReader(htmlR)
		article.Description = ""
		if errdoc == nil {
			textRune := []rune(strings.TrimSpace(doc.Text()))
			if len(textRune) > 150 {
				article.Description = string(textRune[:150])
			} else {
				article.Description = string(textRune)
			}
		}
	}

	if category != nil {
		article.CategoryId = category.Id
	}

	if err = article.Save(config.DB); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "发布成功",
		"data": article,
	})

}
