package main

import (
	"fmt"
<<<<<<< HEAD
<<<<<<< HEAD
	"log"
=======
>>>>>>> cdc3d9a... es5 support completed
=======
	"log"
>>>>>>> e15b2fb... improve
	"strings"
	"syscall/js"
)

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
const src = `function hoge() {console.log(window)}
`

// Parser ...
type Parser struct {
	stack []js.Value
=======
// Parser ...
type Parser struct {
	stack []js.Value
	lines []string
>>>>>>> cdc3d9a... es5 support completed
=======
const src = `
console.log(window);
=======
const src = `function hoge() {console.log(window)}
>>>>>>> ce48c93... improve
`

// Parser ...
type Parser struct {
	stack []js.Value
>>>>>>> e15b2fb... improve
	err   error
}

// ParseProgram ...
func (p *Parser) ParseProgram(obj js.Value) ([]string, error) {
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e15b2fb... improve
	res := p.parseArray(obj.Get("body"))
	if p.err != nil {
		return nil, p.err
	}
	return res, nil
<<<<<<< HEAD
=======
	p.parseArray(obj.Get("body"))
	return nil, p.err
>>>>>>> cdc3d9a... es5 support completed
=======
>>>>>>> e15b2fb... improve
}

func (p *Parser) push(obj js.Value) {
	p.stack = append(p.stack, obj)
}

func (p *Parser) pop() js.Value {
	popped := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return popped
}

func (p *Parser) indent() string {
	return strings.Repeat("  ", len(p.stack))
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e15b2fb... improve
func (p *Parser) parseLiteral(obj js.Value) []string {
	console.Call("log", p.indent(), "Literal:", obj)
	res := "nil"
	switch v := obj.Get("value"); v.Type() {
	case js.TypeNull:
	case js.TypeUndefined:
		res = "js.Undefined()"
	case js.TypeBoolean:
		res = fmt.Sprint(v.Bool())
	case js.TypeNumber:
		res = fmt.Sprint(v.Float())
	case js.TypeString:
		res = fmt.Sprintf("%q", v.String())
	case js.TypeObject:
		fallthrough
	default:
		p.err = fmt.Errorf("unsupported literal: %v", v.Type())
	}
	return []string{res}
}

<<<<<<< HEAD
<<<<<<< HEAD
func (p *Parser) parseVariableDeclaration(obj js.Value) []string {
	console.Call("log", p.indent(), "VariableDeclaration:", obj)
	p.push(obj)
	defer p.pop()
	for i, v := range p.stack {
		log.Print(i, v.Get("type").String())
	}
	kind := obj.Get("kind").String()
	decls := obj.Get("declarations")
	if kind == "let" {
		if len(p.stack) > 1 {
			// := statement
			res := []string{}
			for _, v := range p.parseArray(decls) {
				res = append(res, strings.Replace(v, "=", ":=", 1))
			}
			return res
		}
		kind = "var"
	}
	if decls.Length() == 1 {
		res := p.parseArray(decls)
		return append([]string{fmt.Sprintf("%s %s", kind, res[0])}, res[1:]...)
	}
	res := []string{kind + " ("}
	res = append(res, p.parseArray(decls)...)
	res = append(res, ")")
	return res
}

func (p *Parser) parseVariableDeclarator(obj js.Value) []string {
	console.Call("log", p.indent(), "VariableDeclarator:", obj)
	p.push(obj)
	defer p.pop()
	id := p.parseIdentifier(obj.Get("id"))[0]
	init := obj.Get("init")
	if init.IsNull() {
		return []string{fmt.Sprintf("%s = nil", id)}
	}
	res := p.parseStatement(init)
	return append([]string{fmt.Sprintf("%s = %s", id, res[0])}, res[1:]...)
}

func (p *Parser) parseIdentifier(obj js.Value) []string {
	console.Call("log", p.indent(), "Identifier:", obj)
	name := obj.Get("name").String()
	if name == "window" {
		name = "js.Global()"
	}
	return []string{name}
}

func (p *Parser) parseMemberExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "MemberExpression:", obj)
	res := p.parseStatement(obj.Get("object"))
	res = append(res, p.parseStatement(obj.Get("property"))...)
	return res
}

func buildGetChain(args []string) string {
	res := []string{"js.Global()"}
	for _, v := range args {
		res = append(res, fmt.Sprintf("Get(%q)", v))
	}
	return strings.Join(res, ".")
}

func (p *Parser) parseCallExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "CallExpression:", obj)
	p.push(obj)
	defer p.pop()
	callee := obj.Get("callee")
	args := p.parseArray(obj.Get("arguments"))

	switch callee.Get("type").String() {
	default:
		res := p.parseStatement(callee)
		res[len(res)-1] += "()"
		return res
	case "MemberExpression":
		static := p.parseStatement(callee)
		parent, fn := static[:len(static)-1], static[len(static)-1]
		log.Println(parent, fn)
		switch len(args) {
		case 0:
			return []string{
				buildGetChain(parent) + fmt.Sprintf(".Call(%q)", fn),
			}
		case 1:
			return []string{
				buildGetChain(parent) + fmt.Sprintf(".Call(%q, %s)", fn, strings.Join(args, ", ")),
			}
		default:
		}
		res := []string{buildGetChain(parent) + fmt.Sprintf(".Call(%q,", fn)}
		for _, arg := range args {
			res = append(res, arg+",")
		}
		res = append(res, ")")
		return res
	}
}

func (p *Parser) parseFunctionDeclaration(obj js.Value) []string {
	console.Call("log", p.indent(), "FunctionDeclaration:", obj)
	p.push(obj)
	defer p.pop()
	params := p.parseArray(obj.Get("params"))
	res := []string{fmt.Sprintf("func %s(%s) {",
		p.parseIdentifier(obj.Get("id"))[0],
		strings.Join(params, ", "),
	)}
	res = append(res, p.parseStatement(obj.Get("body"))...)
	res = append(res, "}")
	return res
}

func (p *Parser) parseFunctionExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "FunctionExpression:", obj)
	p.push(obj)
	defer p.pop()
	res := []string{"func(){"}
	res = append(res, p.parseStatement(obj.Get("body"))...)
	res = append(res, "}")
	return res
}

func (p *Parser) parseArrowFunctionExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "ArrowFunctionExpression:", obj)
	p.push(obj)
	defer p.pop()
	res := []string{"func(){"}
	res = append(res, p.parseStatement(obj.Get("body"))...)
	res = append(res, "}")
	return res
=======
=======
>>>>>>> e15b2fb... improve
func (p *Parser) parseVariableDeclaration(obj js.Value) {
=======
func (p *Parser) parseVariableDeclaration(obj js.Value) []string {
>>>>>>> 686d81a... wip
	console.Call("log", p.indent(), "VariableDeclaration:", obj)
	p.push(obj)
	defer p.pop()
	for i, v := range p.stack {
		log.Print(i, v.Get("type").String())
	}
	kind := obj.Get("kind").String()
	decls := obj.Get("declarations")
	if kind == "let" {
		if len(p.stack) > 1 {
			// := statement
			res := []string{}
			for _, v := range p.parseArray(decls) {
				res = append(res, strings.Replace(v, "=", ":=", 1))
			}
			return res
		}
		kind = "var"
	}
	if decls.Length() == 1 {
		res := p.parseArray(decls)
		return append([]string{fmt.Sprintf("%s %s", kind, res[0])}, res[1:]...)
	}
	res := []string{kind + " ("}
	res = append(res, p.parseArray(decls)...)
	res = append(res, ")")
	return res
}

func (p *Parser) parseVariableDeclarator(obj js.Value) []string {
	console.Call("log", p.indent(), "VariableDeclarator:", obj)
	p.push(obj)
	defer p.pop()
	id := p.parseIdentifier(obj.Get("id"))[0]
	init := obj.Get("init")
	if init.IsNull() {
		return []string{fmt.Sprintf("%s = nil", id)}
	}
	res := p.parseStatement(init)
	return append([]string{fmt.Sprintf("%s = %s", id, res[0])}, res[1:]...)
}

func (p *Parser) parseIdentifier(obj js.Value) []string {
	console.Call("log", p.indent(), "Identifier:", obj)
	name := obj.Get("name").String()
	if name == "window" {
		name = "js.Global()"
	}
	return []string{name}
}

func (p *Parser) parseMemberExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "MemberExpression:", obj)
	res := p.parseStatement(obj.Get("object"))
	res = append(res, p.parseStatement(obj.Get("property"))...)
	return res
}

func buildGetChain(args []string) string {
	res := []string{"js.Global()"}
	for _, v := range args {
		res = append(res, fmt.Sprintf("Get(%q)", v))
	}
	return strings.Join(res, ".")
}

func (p *Parser) parseCallExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "CallExpression:", obj)
	p.push(obj)
	defer p.pop()
	callee := obj.Get("callee")
	args := p.parseArray(obj.Get("arguments"))

	switch callee.Get("type").String() {
	default:
		res := p.parseStatement(callee)
		res[len(res)-1] += "()"
		return res
	case "MemberExpression":
		static := p.parseStatement(callee)
		parent, fn := static[:len(static)-1], static[len(static)-1]
		log.Println(parent, fn)
		switch len(args) {
		case 0:
			return []string{
				buildGetChain(parent) + fmt.Sprintf(".Call(%q)", fn),
			}
		case 1:
			return []string{
				buildGetChain(parent) + fmt.Sprintf(".Call(%q, %s)", fn, strings.Join(args, ", ")),
			}
		default:
		}
		res := []string{buildGetChain(parent) + fmt.Sprintf(".Call(%q,", fn)}
		for _, arg := range args {
			res = append(res, arg+",")
		}
		res = append(res, ")")
		return res
	}
}

func (p *Parser) parseFunctionDeclaration(obj js.Value) []string {
	console.Call("log", p.indent(), "FunctionDeclaration:", obj)
	p.push(obj)
	defer p.pop()
	params := p.parseArray(obj.Get("params"))
	res := []string{fmt.Sprintf("func %s(%s) {",
		p.parseIdentifier(obj.Get("id"))[0],
		strings.Join(params, ", "),
	)}
	res = append(res, p.parseStatement(obj.Get("body"))...)
	res = append(res, "}")
	return res
}

func (p *Parser) parseFunctionExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "FunctionExpression:", obj)
	p.push(obj)
	defer p.pop()
<<<<<<< HEAD
	p.parseStatement(obj.Get("body"))
>>>>>>> cdc3d9a... es5 support completed
=======
	res := []string{"func(){"}
	res = append(res, p.parseStatement(obj.Get("body"))...)
	res = append(res, "}")
	return res
}

func (p *Parser) parseArrowFunctionExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "ArrowFunctionExpression:", obj)
	p.push(obj)
	defer p.pop()
	res := []string{"func(){"}
	res = append(res, p.parseStatement(obj.Get("body"))...)
	res = append(res, "}")
	return res
>>>>>>> 686d81a... wip
}

func (p *Parser) parseMethodDefinition(obj js.Value) {
	console.Call("log", p.indent(), "MethodDefinition:", obj)
	p.push(obj)
	defer p.pop()
	p.parseFunctionDeclaration(obj.Get("value"))
}

func (p *Parser) parseSwitchStatement(obj js.Value) {
	console.Call("log", p.indent(), "SwitchStatement:", obj)
	p.push(obj)
	defer p.pop()
	p.parseArray(obj.Get("cases"))
}

func (p *Parser) parseWhileStatement(obj js.Value) {
	console.Call("log", p.indent(), "WhileStatement:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("body"))
}

func (p *Parser) parseDoWhileStatement(obj js.Value) {
	console.Call("log", p.indent(), "DoWhileStatement:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("body"))
}

func (p *Parser) parseForStatement(obj js.Value) {
	console.Call("log", p.indent(), "ForStatement:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("body"))
}

func (p *Parser) parseForInStatement(obj js.Value) {
	console.Call("log", p.indent(), "ForInStatement:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("body"))
}

func (p *Parser) parseIfStatement(obj js.Value) {
	console.Call("log", p.indent(), "IfStatement:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("consequent"))
}

func (p *Parser) parseTryStatement(obj js.Value) {
	console.Call("log", p.indent(), "TryStatement:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("block"))
}

func (p *Parser) parseAssignmentExpression(obj js.Value) {
	console.Call("log", p.indent(), "AssignmentExpression:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("left"))
	p.parseStatement(obj.Get("right"))
}

func (p *Parser) parseClassDeclaration(obj js.Value) {
	console.Call("log", p.indent(), "ClassDeclaration:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("body"))
}

<<<<<<< HEAD
<<<<<<< HEAD
func (p *Parser) parseAwaitExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "AwaitExpression:", obj)
	p.push(obj)
	defer p.pop()
	res := p.parseCallExpression(obj.Get("argument"))
	return []string{fmt.Sprintf("jsutil.Await(%s)", strings.Join(res, ", "))}
=======
func (p *Parser) parseAwaitExpression(obj js.Value) {
	console.Call("log", p.indent(), "AwaitExpression:", obj)
	p.push(obj)
	defer p.pop()
	p.parseCallExpression(obj.Get("argument"))
>>>>>>> cdc3d9a... es5 support completed
=======
func (p *Parser) parseAwaitExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "AwaitExpression:", obj)
	p.push(obj)
	defer p.pop()
	res := p.parseCallExpression(obj.Get("argument"))
	return []string{fmt.Sprintf("jsutil.Await(%s)", strings.Join(res, ", "))}
>>>>>>> e15b2fb... improve
}

func (p *Parser) parseReturnStatement(obj js.Value) {
	console.Call("log", p.indent(), "ReturnStatement:", obj)
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e15b2fb... improve
func (p *Parser) parseArray(body js.Value, suffix ...string) []string {
	res := []string{}
	for i := 0; i < body.Length(); i++ {
		res = append(res,
			append(p.parseStatement(body.Index(i)), suffix...)...)
	}
	return res
}

func (p *Parser) parseStatement(obj js.Value) []string {
	res := []string{}
<<<<<<< HEAD
=======
func (p *Parser) parseStatement(obj js.Value) {
>>>>>>> cdc3d9a... es5 support completed
=======
>>>>>>> e15b2fb... improve
	switch v := obj.Get("type").String(); v {
	default:
		console.Call("log", "unknown expression type:", obj)
		p.err = fmt.Errorf("unknown expression type: %s", v)
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e15b2fb... improve
	case "EmptyStatement":
	case "Literal":
		res = append(res, p.parseLiteral(obj)...)
	case "Identifier":
		res = append(res, p.parseIdentifier(obj)...)
<<<<<<< HEAD
	case "BlockStatement":
		console.Call("log", p.indent(), "BlockStatement:", obj)
		res = append(res, p.parseArray(obj.Get("body"))...)
	case "ClassBody":
		console.Call("log", p.indent(), "ClassBody:", obj)
		res = append(res, p.parseArray(obj.Get("body"))...)
	case "ExpressionStatement":
		console.Call("log", p.indent(), "ExpressionStatement:", obj)
		res = append(res, p.parseStatement(obj.Get("expression"))...)
	case "VariableDeclarator":
		res = append(res, p.parseVariableDeclarator(obj)...)
	case "VariableDeclaration":
		res = append(res, p.parseVariableDeclaration(obj)...)
	case "FunctionExpression":
		res = append(res, p.parseFunctionExpression(obj)...)
=======
=======
>>>>>>> e15b2fb... improve
	case "BlockStatement":
		console.Call("log", p.indent(), "BlockStatement:", obj)
		res = append(res, p.parseArray(obj.Get("body"))...)
	case "ClassBody":
		console.Call("log", p.indent(), "ClassBody:", obj)
		res = append(res, p.parseArray(obj.Get("body"))...)
	case "ExpressionStatement":
		console.Call("log", p.indent(), "ExpressionStatement:", obj)
		res = append(res, p.parseStatement(obj.Get("expression"))...)
	case "VariableDeclarator":
		res = append(res, p.parseVariableDeclarator(obj)...)
	case "VariableDeclaration":
		res = append(res, p.parseVariableDeclaration(obj)...)
	case "FunctionExpression":
<<<<<<< HEAD
<<<<<<< HEAD
		p.parseStatement(obj.Get("body"))
>>>>>>> cdc3d9a... es5 support completed
=======
		console.Call("log", p.indent(), "FunctionExpression:", obj)
		res = append(res, p.parseStatement(obj.Get("body"))...)
>>>>>>> e15b2fb... improve
=======
		res = append(res, p.parseFunctionExpression(obj)...)
>>>>>>> 686d81a... wip
	case "ThisExpression":
		console.Call("log", p.indent(), "ThisExpression:", obj)
	case "ObjectExpression":
		console.Call("log", p.indent(), "ObjectExpression:", obj)
	case "BinaryExpression":
		console.Call("log", p.indent(), "BinaryExpression:", obj)
	case "LogicalExpression":
		console.Call("log", p.indent(), "LogicalExpression:", obj)
	case "ArrayExpression":
		console.Call("log", p.indent(), "ArrayExpression:", obj)
<<<<<<< HEAD
<<<<<<< HEAD
=======
	case "Identifier":
		console.Call("log", p.indent(), "Identifier:", obj)
>>>>>>> cdc3d9a... es5 support completed
=======
>>>>>>> e15b2fb... improve
	case "ThrowStatement":
		console.Call("log", p.indent(), "ThrowStatement:", obj)
	case "ConditionalExpression":
		console.Call("log", p.indent(), "ConditionalExpression:", obj)
	case "CallExpression":
<<<<<<< HEAD
<<<<<<< HEAD
		res = append(res, p.parseCallExpression(obj)...)
		log.Println(res)
=======
		p.parseCallExpression(obj)
>>>>>>> cdc3d9a... es5 support completed
=======
		res = append(res, p.parseCallExpression(obj)...)
		log.Println(res)
>>>>>>> e15b2fb... improve
	case "AssignmentExpression":
		p.parseAssignmentExpression(obj)
	case "AwaitExpression":
		p.parseAwaitExpression(obj)
	case "ArrowFunctionExpression":
		p.parseArrowFunctionExpression(obj)
<<<<<<< HEAD
<<<<<<< HEAD
=======
	case "EmptyStatement":
	case "Literal":
		console.Call("log", p.indent(), "Literal:", obj)
>>>>>>> cdc3d9a... es5 support completed
=======
>>>>>>> e15b2fb... improve
	case "SequenceExpression":
		console.Call("log", p.indent(), "SequenceExpression:", obj)
	case "UpdateExpression":
		console.Call("log", p.indent(), "UpdateExpression:", obj)
	case "ContinueStatement":
		console.Call("log", p.indent(), "ContinueStatement:", obj)
	case "BreakStatement":
		console.Call("log", p.indent(), "BreakStatement:", obj)
	case "NewExpression":
		console.Call("log", p.indent(), "NewExpression:", obj)
	case "ComputedMemberExpression":
		console.Call("log", p.indent(), "ComputedMemberExpression:", obj)
	case "MemberExpression":
<<<<<<< HEAD
<<<<<<< HEAD
		res = append(res, p.parseMemberExpression(obj)...)
	case "ReturnStatement":
		p.parseReturnStatement(obj)
<<<<<<< HEAD
	case "FunctionDeclaration":
		res = append(res, p.parseFunctionDeclaration(obj)...)
=======
		console.Call("log", p.indent(), "MemberExpression:", obj)
=======
		res = append(res, p.parseMemberExpression(obj)...)
>>>>>>> e15b2fb... improve
	case "ReturnStatement":
		p.parseReturnStatement(obj)
	case "VariableDeclaration":
		p.parseVariableDeclaration(obj)
=======
>>>>>>> 686d81a... wip
	case "FunctionDeclaration":
<<<<<<< HEAD
		p.parseFunctionDeclaration(obj)
>>>>>>> cdc3d9a... es5 support completed
=======
		res = append(res, p.parseFunctionDeclaration(obj)...)
>>>>>>> ce48c93... improve
	case "ClassDeclaration":
		p.parseClassDeclaration(obj)
	case "SwitchStatement":
		p.parseSwitchStatement(obj)
	case "SwitchCase":
		console.Call("log", p.indent(), "SwitchCase:", obj)
	case "WhileStatement":
		p.parseWhileStatement(obj)
	case "DoWhileStatement":
		p.parseDoWhileStatement(obj)
	case "ForStatement":
		p.parseForStatement(obj)
	case "ForInStatement":
		p.parseForInStatement(obj)
	case "IfStatement":
		p.parseIfStatement(obj)
	case "TryStatement":
		p.parseTryStatement(obj)
	case "MethodDefinition":
		p.parseMethodDefinition(obj)
	}
<<<<<<< HEAD
<<<<<<< HEAD
	return res
=======
}

func (p *Parser) parseArray(body js.Value) {
	for i := 0; i < body.Length(); i++ {
		p.parseStatement(body.Index(i))
	}
>>>>>>> cdc3d9a... es5 support completed
=======
	return res
>>>>>>> e15b2fb... improve
}
