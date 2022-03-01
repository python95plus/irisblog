package tags

import (
	"fmt"
	"irisblog/config"
	"irisblog/provider"

	"github.com/iris-contrib/pongo2"
)

type tagCategoryListNode struct {
	name    string
	args    map[string]pongo2.IEvaluator
	wrapper *pongo2.NodeWrapper
}

func (node *tagCategoryListNode) Execute(ctx *pongo2.ExecutionContext, writer pongo2.TemplateWriter) *pongo2.Error {
	if config.DB == nil {
		return nil
	}
	args, err := parseArgs(node.args, ctx)
	if err != nil {
		return err
	}

	parentId := uint(0)
	if args["parentId"] != nil {
		parentId = uint(args["parentId"].Integer())
	}
	categoryList, _ := provider.GetCategories(parentId)

	ctx.Private[node.name] = categoryList
	// ctx.Private是pongo2输出给template的,ctx.Public是由控制模型中ctx.ViewData的ctx iris.Context

	//execute
	node.wrapper.Execute(ctx, writer)
	return nil
}

// categoryList articleCategories with type="1" parentId="0"
func TagCategoryListParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	// doc为categoryList前面的节点,start为categoryList节点, arguments为categories节点
	tagNode := &tagCategoryListNode{
		args: make(map[string]pongo2.IEvaluator),
	}

	nameToken := arguments.MatchType(pongo2.TokenIdentifier)
	//{% categoryList categories with type="2" parentId="0" %}
	// arguments.MatchType执行后,idx会自动加1,本来在categories加1是with节点,而在加1之前已经把值返回给nameToken(实际是给一个临时变量)

	if nameToken == nil {
		return nil, arguments.Error("categoryList-tag needs a accept name.", nil)
	}

	tagNode.name = nameToken.Val
	args, err := parseWith(arguments)
	if err != nil {
		return nil, err
	}

	tagNode.args = args

	for arguments.Remaining() > 0 { //节点数量减去
		return nil, arguments.Error("Malformed categoryList-tag arguments.", nil)
	}
	wrapper, endtagargs, err := doc.WrapUntilTag("endcategoryList")
	if err != nil {
		return nil, err
	}

	if endtagargs.Remaining() > 0 {
		endtagnameToken := endtagargs.MatchType(pongo2.TokenIdentifier)
		if endtagnameToken != nil {
			if endtagnameToken.Val != nameToken.Val {
				return nil, endtagargs.Error(fmt.Sprintf("Name for 'endcategoryList' must equal to 'categoryList'-tag's name ('%s' != '%s').",
					nameToken.Val, endtagnameToken.Val), nil)
			}
		}

		if endtagnameToken == nil || endtagargs.Remaining() > 0 {
			return nil, endtagargs.Error("Either no or only one argument (identifier) allowed for 'endcategoryList'.", nil)
		}
	}

	tagNode.wrapper = wrapper
	return tagNode, nil
}
