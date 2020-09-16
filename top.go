package main

import (
	"strings"
	"syscall/js"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/jsutil"
)

//go:generate spago generate -c Top -p main top.html

type (
	// Object ...
	Object = map[string]interface{}
	// Array ...
	Array = []interface{}
)

// Top  ...
type Top struct {
	spago.Core
	JsCode string
}

const src = `
const name = (hoge) => {
    console.log(hoge);
};
`

// OnSubmit ...
func (c *Top) OnSubmit(ev js.Value) {
	ev.Call("preventDefault")
	params := jsutil.Form2Go(ev.Get("target"))
	tree := esprima.Call("parseScript", params["js"])
	console.Call("log", tree)
	(&Parser{}).parseProgram(tree)
}

// Parser ...
type Parser struct {
	indent int
}

// Indent ...
func (p *Parser) Indent() string {
	return strings.Repeat("  ", p.indent)
}

func (p *Parser) parseExpressionStatement(obj js.Value) {
	p.indent++
	defer func() { p.indent-- }()
	switch obj.Get("type").String() {
	default:
		console.Call("log", p.Indent(), "unknown type:", obj)
	case "CallExpression":
		p.parseCallExpression(obj)
	case "AssignmentExpression":
		p.parseAssignmentExpression(obj)
	case "AwaitExpression":
		p.parseAwaitExpression(obj)
	case "ArrowFunctionExpression":
		p.parseArrowFunctionExpression(obj)
	}
}

func (p *Parser) parseVariableDeclaration(obj js.Value) {
	switch obj.Get("type").String() {
	default:
		console.Call("log", p.Indent(), "unknown type:", obj)
	case "VariableDeclaration":
		kind := obj.Get("kind").String()
		declarations := obj.Get("declarations")
		for i := 0; i < declarations.Length(); i++ {
			decl := declarations.Index(i)
			switch decl.Get("type").String() {
			default:
				console.Call("log", p.Indent(), "unknown type:", decl)
			case "VariableDeclaration":
				console.Call("log", p.Indent(), "VariableDeclaration:", kind, decl)
			case "VariableDeclarator":
				console.Call("log", p.Indent(), "VariableDeclarator:", kind, decl)
				p.parseExpressionStatement(decl.Get("init"))
			}
		}
	}
}
func (p *Parser) parseCallExpression(obj js.Value) {
	console.Call("log", p.Indent(), "CallExpression:", obj)
}

func (p *Parser) parseFunctionDeclaration(obj js.Value) {
	console.Call("log", p.Indent(), "FunctionDeclaration:", obj)
	p.parseBody(obj.Get("body"))
}

func (p *Parser) parseArrowFunctionExpression(obj js.Value) {
	console.Call("log", p.Indent(), "FunctionDeclaration:", obj)
	p.parseBody(obj.Get("body"))
}

func (p *Parser) parseMethodDefinition(obj js.Value) {
	console.Call("log", p.Indent(), "MethodDefinition:", obj)
	p.parseFunctionDeclaration(obj.Get("value"))
}

func (p *Parser) parseWhileStatement(obj js.Value) {
	console.Call("log", p.Indent(), "WhileStatement:", obj)
	p.parseBody(obj.Get("body"))
}

func (p *Parser) parseForStatement(obj js.Value) {
	console.Call("log", p.Indent(), "ForStatement:", obj)
	p.parseBody(obj.Get("body"))
}

func (p *Parser) parseIfStatement(obj js.Value) {
	console.Call("log", p.Indent(), "IfStatement:", obj)
	p.parseBody(obj.Get("consequent"))
}

func (p *Parser) parseTryStatement(obj js.Value) {
	console.Call("log", p.Indent(), "TryStatement:", obj)
	p.parseBody(obj.Get("block"))
}

func (p *Parser) parseAssignmentExpression(obj js.Value) {
	console.Call("log", p.Indent(), "AssignmentExpression:", obj)
}

func (p *Parser) parseClassDeclaration(obj js.Value) {
	console.Call("log", p.Indent(), "ClassDeclaration:", obj)
	p.parseBody(obj.Get("body"))
}

func (p *Parser) parseAwaitExpression(obj js.Value) {
	console.Call("log", p.Indent(), "AwaitExpression:", obj)
	p.indent++
	defer func() { p.indent-- }()
	p.parseCallExpression(obj.Get("argument"))
}

func (p *Parser) parseBody(obj js.Value) {
	p.indent++
	defer func() { p.indent-- }()
	switch obj.Get("type").String() {
	default:
		console.Call("log", p.Indent(), "unknown type:", obj)
	case "ExpressionStatement":
		p.parseExpressionStatement(obj.Get("expression"))
	case "BlockStatement":
		p.parseBodyArray(obj.Get("body"))
	case "ClassBody":
		p.parseBodyArray(obj.Get("body"))
	}
}

func (p *Parser) parseBodyArray(body js.Value) {
	for i := 0; i < body.Length(); i++ {
		statement := body.Index(i)
		switch statement.Get("type").String() {
		default:
			console.Call("log", p.Indent(), "unknown type:", statement)
		case "ExpressionStatement":
			p.parseExpressionStatement(statement.Get("expression"))
		case "VariableDeclaration":
			p.parseVariableDeclaration(statement)
		case "FunctionDeclaration":
			p.parseFunctionDeclaration(statement)
		case "ClassDeclaration":
			p.parseClassDeclaration(statement)
		case "WhileStatement":
			p.parseWhileStatement(statement)
		case "ForStatement":
			p.parseForStatement(statement)
		case "IfStatement":
			p.parseIfStatement(statement)
		case "TryStatement":
			p.parseTryStatement(statement)
		case "MethodDefinition":
			p.parseMethodDefinition(statement)
		}
	}
}

func (p *Parser) parseProgram(obj js.Value) ([]string, error) {
	p.parseBodyArray(obj.Get("body"))
	return nil, nil
}
