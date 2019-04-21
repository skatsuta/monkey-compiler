package eval

import (
	"strconv"

	"github.com/skatsuta/monkey-compiler/ast"
	"github.com/skatsuta/monkey-compiler/object"
	"github.com/skatsuta/monkey-compiler/token"
)

const (
	// FuncNameQuote is a name for quote function.
	FuncNameQuote = "quote"
	// FuncNameUnquote is a name for unquote function.
	FuncNameUnquote = "unquote"
)

func quote(node ast.Node, env object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(quoted ast.Node, env object.Environment) ast.Node {
	modifier := func(node ast.Node) ast.Node {
		call, ok := node.(*ast.CallExpression)
		if !ok || call.Function.TokenLiteral() != FuncNameUnquote || len(call.Arguments) != 1 {
			return node
		}

		unquoted := Eval(call.Arguments[0], env)
		return convertObjectToASTNode(unquoted)
	}
	return ast.Modify(quoted, modifier)
}

func convertObjectToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		base := 10
		t := token.Token{
			Type:    token.Int,
			Literal: strconv.FormatInt(obj.Value, base),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}
	case *object.Boolean:
		var t token.Token
		if obj.Value {
			t = token.Token{Type: token.True, Literal: "true"}
		} else {
			t = token.Token{Type: token.False, Literal: "false"}
		}
		return &ast.Boolean{Token: t, Value: obj.Value}
	case *object.Quote:
		return obj.Node
	default:
		return nil
	}
}
