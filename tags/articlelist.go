package tags

import (
	"fmt"
	"irisblog/config"
	"irisblog/model"
	"irisblog/provider"

	"github.com/iris-contrib/pongo2"
)

type tagArticleListNode struct {
	name   string
	args   map[string]pongo2.IEvaluator
	wraper *pongo2.NodeWrapper
}

func (node *tagArticleListNode) Execute(ctx *pongo2.ExecutionContext, writer pongo2.TemplateWriter) *pongo2.Error {
	// {% articleList popularArticles with type="list" order="views desc" limit="6" %}
	if config.DB == nil {
		return nil
	}

	args, err := parseArgs(node.args, ctx)
	if err != nil {
		return err
	}

	order := "id desc"
	if args["order"] != nil {
		order = args["order"].String()
	}

	limit := 10
	if args["limit"] != nil {
		limit = args["limit"].Integer()
		if limit > 100 {
			limit = 100
		} else if limit < 1 {
			limit = 1
		}
	}

	listType := "list"
	if args["type"] != nil {
		listType = args["type"].String()
	}

	var articleList []*model.Article
	categoryId := uint(0)
	if listType == "related" {
		artilceDetail, ok := ctx.Public["article"].(*model.Article)
		if ok {
			categoryId = artilceDetail.CategoryId
		}

		articleList, _ = provider.GetArticleList(categoryId, order, 0, limit)
	} else {
		// categoryId uint, order string, currentPage int, pageSize int
		articleList, _ = provider.GetArticleList(categoryId, order, 0, limit)
	}

	ctx.Private[node.name] = articleList

	node.wraper.Execute(ctx, writer)
	return nil
}

func TagArticleListParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {

	tagNode := &tagArticleListNode{
		args: make(map[string]pongo2.IEvaluator),
	}

	nameToken := arguments.MatchType(pongo2.TokenIdentifier)
	if nameToken == nil {
		return nil, arguments.Error("articleList-tag needs a accept name.", nil)
	}
	tagNode.name = nameToken.Val

	args, err := parseWith(arguments)
	if err != nil {
		return nil, err
	}
	tagNode.args = args

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed articleList-tag arguments.", nil)
	}

	wraper, endtagargs, err := doc.WrapUntilTag("endarticleList")
	if err != nil {
		return nil, err
	}

	if endtagargs.Remaining() > 0 {
		endtagnameToken := endtagargs.MatchType(pongo2.TokenIdentifier)
		if endtagnameToken != nil {
			if endtagnameToken.Val != nameToken.Val {
				return nil, endtagargs.Error(fmt.Sprintf("Name for 'endarticleList' must equal to 'articleList'-tag's name ('%s' != '%s').",
					nameToken.Val, endtagnameToken.Val), nil)
			}
		}

		if endtagnameToken == nil || endtagargs.Remaining() > 0 {
			return nil, endtagargs.Error("Either no or only one argument (identifier) allowed for 'endarticleList'.", nil)
		}
	}

	tagNode.wraper = wraper

	return tagNode, nil
}
