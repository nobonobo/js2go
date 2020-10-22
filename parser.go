package main

import (
	"fmt"
	"log"
	"strings"
	"syscall/js"
)

const src = `async function on_click() {
	let device = await navigator.bluetooth.requestDevice({ filters: [{ services: ["battery_service"] }] })
}`

// Parser ...
type Parser struct {
	stack []stack
	err   error
}

type stack struct {
	obj    js.Value
	scope  map[string]bool
	define bool
	prefix []string
	suffix []string
}

func (s *stack) append(src []string) []string {
	res := s.prefix
	res = append(res, src...)
	res = append(res, s.suffix...)
	return res
}

// ParseProgram ...
func (p *Parser) ParseProgram(obj js.Value) ([]string, error) {
	p.push(obj)
	defer p.pop()
	res := p.parseArray(obj.Get("body"))
	if p.err != nil {
		return nil, p.err
	}
	return res, nil
}

func (p *Parser) push(obj js.Value) {
	p.stack = append(p.stack, stack{obj: obj, scope: map[string]bool{}})
}

func (p *Parser) pop() stack {
	popped := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return popped
}

func (p *Parser) indent() string {
	return strings.Repeat("  ", len(p.stack))
}

func (p *Parser) setDefineMode(b bool) {
	p.stack[len(p.stack)-1].define = b
}

func (p *Parser) getDefineMode() bool {
	return p.stack[len(p.stack)-1].define
}

func (p *Parser) define(sym string, native bool) {
	p.stack[len(p.stack)-1].scope[sym] = native
}

func (p *Parser) defined(sym string) bool {
	for i := len(p.stack) - 1; i >= 0; i-- {
		_, ok := p.stack[i].scope[sym]
		if ok {
			return true
		}
	}
	return false
}

func (p *Parser) definedAndType(sym string) (bool, bool) {
	for i := len(p.stack) - 1; i >= 0; i-- {
		tp, ok := p.stack[i].scope[sym]
		if ok {
			return true, tp
		}
	}
	return false, false
}

func (p *Parser) parseIdentifier(obj js.Value) string {
	return obj.Get("name").String()
}

func (p *Parser) parseIdentifierRef(obj js.Value) []string {
	console.Call("log", p.indent(), "Identifier:", obj)
	name := obj.Get("name").String()
	if name == "window" {
		return []string{"js.Global()"}
	}
	if p.defined(name) {
		return []string{name}
	}
	return []string{fmt.Sprintf("js.Global().Get(%q)", name)}
}

func (p *Parser) parseProperty(obj js.Value) []string {
	console.Call("log", p.indent(), "Property:", obj)
	key := p.parseStatement(obj.Get("key"))[0]
	value := p.parseStatement(obj.Get("value"))[0]
	res := []string{fmt.Sprintf("%q: %s", key, value)}
	return res
}

func (p *Parser) parseArrayExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "ArrayExpression:", obj)
	elements := p.parseArray(obj.Get("elements"))
	res := []string{"[]interface{}{"}
	switch len(elements) {
	case 0:
		res[0] += "}"
	case 1:
		res[0] += elements[0] + "}"
	default:
		for _, v := range elements {
			res = append(res, v+",")
		}
		res = append(res, "}")
	}
	return res
}

func (p *Parser) parseThisExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "ThisExpression:", obj)
	res := []string{}
	return res
}

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

func (p *Parser) parseVariableDeclaration(obj js.Value) []string {
	console.Call("log", p.indent(), "VariableDeclaration:", obj)
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
	id := p.parseIdentifier(obj.Get("id"))
	p.define(id, true)
	init := obj.Get("init")
	if init.IsNull() {
		return []string{fmt.Sprintf("%s = nil", id)}
	}
	switch init.Get("type").String() {
	case "AwaitExpression":
		id += ", err"
	case "MemberExpression":
		prop := p.parseMemberExpression(init)
		return []string{fmt.Sprintf("%s = %s.Get(%q)", id, prop[0], prop[1])}
	default:
	}
	res := p.parseStatement(init)
	return append([]string{fmt.Sprintf("%s = %s", id, res[0])}, res[1:]...)
}

func (p *Parser) buildGetChain(args []string) string {
	res := []string{}
	if p.stack[len(p.stack)-1].obj.Get("type").String() == "CallExpression" {
		res = append(res, "js.Global()")
	} else {
		res = append(res, args[0])
		args = args[1:]
	}
	for _, v := range args {
		res = append(res, fmt.Sprintf("Get(%q)", v))
	}
	return strings.Join(res, ".")
}

func (p *Parser) parseMemberExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "MemberExpression:", obj)
	target := obj.Get("object")
	res := []string{}
	switch target.Get("type").String() {
	case "MemberExpression":
		res = append(res, p.buildGetChain(p.parseStatement(target)))
	default:
		res = append(res, p.parseStatement(target)...)
	}
	res = append(res, p.parseIdentifier(obj.Get("property")))
	return res
}

func (p *Parser) parseCallExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "CallExpression:", obj)
	callee := obj.Get("callee")
	args := obj.Get("arguments")
	switch callee.Get("type").String() {
	default:
		res := p.parseStatement(callee)
		switch args.Length() {
		case 0:
			res[len(res)-1] += "()"
		case 1:
			param := p.parseStatement(args.Index(0))
			res[len(res)-1] += fmt.Sprintf("(%s, ", param[0])
			res = append(res, param[1:]...)
			res[len(res)-1] += ")"
		default:
			res[len(res)-1] += "("
			for i := 0; i < args.Length(); i++ {
				res = append(res, p.parseStatement(args.Index(i))...)
				res[len(res)-1] += ","
			}
			res = append(res, ")")
		}
		return res
	case "Identifier":
		sym := p.parseIdentifier(callee)
		res := []string{}
		if def, tp := p.definedAndType(sym); def {
			if tp {
				res = append(res, fmt.Sprintf("%s(", sym))
			} else {
				res = append(res, fmt.Sprintf("%s.Invoke(", sym))
			}
			switch args.Length() {
			case 0:
				res[len(res)-1] += ")"
			case 1:
				res[len(res)-1] += p.parseStatement(args.Index(0))[0]
				res[len(res)-1] += ")"
			default:
				for i := 0; i < args.Length(); i++ {
					res = append(res, p.parseStatement(args.Index(i))...)
					res[len(res)-1] += ","
				}
				res = append(res, ")")
			}
		} else {
			res = append(res, fmt.Sprintf("js.Global().Call(%q", sym))
			log.Println("sym", sym, "args", args, res)
			switch args.Length() {
			case 0:
				res[len(res)-1] += ")"
			case 1:
				res[len(res)-1] += ", "
				res[len(res)-1] += p.parseStatement(args.Index(0))[0]
				res[len(res)-1] += ")"
			default:
				for i := 0; i < args.Length(); i++ {
					res = append(res, p.parseStatement(args.Index(i))...)
				}
				res = append(res, ")")
			}
		}
		return res
	case "MemberExpression":
		static := p.parseStatement(callee)
		res, fn := static[:len(static)-1], static[len(static)-1]
		res[len(res)-1] += fmt.Sprintf(".Call(%q", fn)
		switch args.Length() {
		case 0:
			res[len(res)-1] += ")"
		case 1:
			res[len(res)-1] += ", "
			res[len(res)-1] += p.parseStatement(args.Index(0))[0]
			res[len(res)-1] += ")"
		default:
			for i := 0; i < args.Length(); i++ {
				res = append(res, p.parseStatement(args.Index(i))...)
			}
			res = append(res, ")")
		}
		return res
	}
}

func (p *Parser) parseObjectExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "ObjectExpression:", obj)
	p.setDefineMode(true)
	defer p.setDefineMode(false)
	properties := p.parseArray(obj.Get("properties"))
	res := []string{}
	switch len(properties) {
	case 0:
		res = append(res, "map[string]interface{}{}")
	case 1:
		res = append(res, fmt.Sprintf("map[string]interface{}{%s}", properties[0]))
	default:
		res = append(res, "map[string]interface{}{")
		for _, v := range properties {
			res = append(res, v+",")
		}
		res = append(res, "}")
	}
	return res
}

func (p *Parser) parseParams(obj js.Value) []string {
	console.Call("log", p.indent(), "ParamsArray:", obj)
	res := []string{}
	for i := 0; i < obj.Length(); i++ {
		id := p.parseIdentifier(obj.Index(i))
		p.define(id, false)
		res = append(res, fmt.Sprintf("%s js.Value", id))
	}
	return res
}

func (p *Parser) parseFunctionDeclaration(obj js.Value) []string {
	console.Call("log", p.indent(), "FunctionDeclaration:", obj)
	id := p.parseIdentifier(obj.Get("id"))
	p.define(id, true)
	p.push(obj)
	defer p.pop()
	params := p.parseParams(obj.Get("params"))
	res := []string{fmt.Sprintf("func %s(%s) {",
		id,
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
	params := p.parseParams(obj.Get("params"))
	res := []string{fmt.Sprintf("func(%s) {", strings.Join(params, ", "))}
	res = append(res, p.parseStatement(obj.Get("body"))...)
	res = append(res, "}")
	return res
}

func (p *Parser) parseArrowFunctionExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "ArrowFunctionExpression:", obj)
	p.push(obj)
	defer p.pop()
	params := p.parseParams(obj.Get("params"))
	res := []string{fmt.Sprintf("func(%s) {", strings.Join(params, ", "))}
	res = append(res, p.parseStatement(obj.Get("body"))...)
	res = append(res, "}")
	return res
}

func (p *Parser) parseMethodDefinition(obj js.Value) []string {
	console.Call("log", p.indent(), "MethodDefinition:", obj)
	p.push(obj)
	defer p.pop()
	p.parseFunctionDeclaration(obj.Get("value"))
	res := []string{}
	return res
}

func (p *Parser) parseSwitchStatement(obj js.Value) []string {
	console.Call("log", p.indent(), "SwitchStatement:", obj)
	p.push(obj)
	defer p.pop()
	p.parseArray(obj.Get("cases"))
	res := []string{}
	return res
}

func (p *Parser) parseSwitchCase(obj js.Value) []string {
	console.Call("log", p.indent(), "SwitchCase:", obj)
	res := []string{}
	return res
}

func (p *Parser) parseWhileStatement(obj js.Value) []string {
	console.Call("log", p.indent(), "WhileStatement:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("body"))
	res := []string{}
	return res
}

func (p *Parser) parseDoWhileStatement(obj js.Value) []string {
	console.Call("log", p.indent(), "DoWhileStatement:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("body"))
	res := []string{}
	return res
}

func (p *Parser) parseForStatement(obj js.Value) []string {
	console.Call("log", p.indent(), "ForStatement:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("body"))
	res := []string{}
	return res
}

func (p *Parser) parseForInStatement(obj js.Value) []string {
	console.Call("log", p.indent(), "ForInStatement:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("body"))
	res := []string{}
	return res
}

func (p *Parser) parseIfStatement(obj js.Value) []string {
	console.Call("log", p.indent(), "IfStatement:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("consequent"))
	res := []string{}
	return res
}

func (p *Parser) parseTryStatement(obj js.Value) []string {
	console.Call("log", p.indent(), "TryStatement:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("block"))
	res := []string{}
	return res
}

func (p *Parser) parseAssignmentExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "AssignmentExpression:", obj)
	res := []string{}
	p.parseStatement(obj.Get("left"))
	p.parseStatement(obj.Get("right"))
	return res
}

func (p *Parser) parseClassDeclaration(obj js.Value) []string {
	console.Call("log", p.indent(), "ClassDeclaration:", obj)
	p.push(obj)
	defer p.pop()
	p.parseStatement(obj.Get("body"))
	res := []string{}
	return res
}

func (p *Parser) parseAwaitExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "AwaitExpression:", obj)
	res := p.parseCallExpression(obj.Get("argument"))
	return []string{fmt.Sprintf("jsutil.Await(%s)", strings.Join(res, ", "))}
}

func (p *Parser) parseConditionalExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "ConditionalExpression:", obj)
	res := []string{}
	return res
}

func (p *Parser) parseReturnStatement(obj js.Value) []string {
	console.Call("log", p.indent(), "ReturnStatement:", obj)
	res := []string{}
	return res
}

func (p *Parser) parseSequenceExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "SequenceExpression:", obj)
	res := []string{}
	return res
}

func (p *Parser) parseThrowStatement(obj js.Value) []string {
	console.Call("log", p.indent(), "ThrowStatement:", obj)
	res := []string{}
	return res
}

func (p *Parser) parseUpdateExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "UpdateExpression:", obj)
	res := []string{}
	return res
}

func (p *Parser) parseContinueStatement(obj js.Value) []string {
	console.Call("log", p.indent(), "ContinueStatement:", obj)
	res := []string{}
	return res
}

func (p *Parser) parseBreakStatement(obj js.Value) []string {
	console.Call("log", p.indent(), "BreakStatement:", obj)
	res := []string{}
	return res
}

func (p *Parser) parseNewExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "NewExpression:", obj)
	res := []string{}
	return res
}

func (p *Parser) parseComputedMemberExpression(obj js.Value) []string {
	console.Call("log", p.indent(), "ComputedMemberExpression:", obj)
	res := []string{}
	return res
}

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
	switch v := obj.Get("type").String(); v {
	default:
		console.Call("log", "unknown expression type:", obj)
		p.err = fmt.Errorf("unknown expression type: %s", v)
	case "EmptyStatement":
	case "Literal":
		res = append(res, p.parseLiteral(obj)...)
	case "Identifier":
		if p.getDefineMode() {
			res = append(res, p.parseIdentifier(obj))
		} else {
			res = append(res, p.parseIdentifierRef(obj)...)
		}
	case "Property":
		res = append(res, p.parseProperty(obj)...)
	case "BlockStatement":
		console.Call("log", p.indent(), "BlockStatement:", obj)
		p.push(obj)
		res = append(res, p.parseArray(obj.Get("body"))...)
		p.pop()
	case "ClassBody":
		console.Call("log", p.indent(), "ClassBody:", obj)
		p.push(obj)
		res = append(res, p.parseArray(obj.Get("body"))...)
		p.pop()
	case "ExpressionStatement":
		console.Call("log", p.indent(), "ExpressionStatement:", obj)
		p.push(obj)
		res = append(res, p.parseStatement(obj.Get("expression"))...)
		p.pop()
	case "VariableDeclarator":
		res = append(res, p.parseVariableDeclarator(obj)...)
	case "VariableDeclaration":
		res = append(res, p.parseVariableDeclaration(obj)...)
	case "FunctionExpression":
		res = append(res, p.parseFunctionExpression(obj)...)
	case "ThisExpression":
		res = append(res, p.parseThisExpression(obj)...)
	case "ObjectExpression":
		res = append(res, p.parseObjectExpression(obj)...)
	case "BinaryExpression":
		console.Call("log", p.indent(), "BinaryExpression:", obj)
	case "LogicalExpression":
		console.Call("log", p.indent(), "LogicalExpression:", obj)
	case "ArrayExpression":
		res = append(res, p.parseArrayExpression(obj)...)
	case "ThrowStatement":
		res = append(res, p.parseThrowStatement(obj)...)
	case "ConditionalExpression":
		res = append(res, p.parseConditionalExpression(obj)...)
	case "CallExpression":
		res = append(res, p.parseCallExpression(obj)...)
	case "AssignmentExpression":
		res = append(res, p.parseAssignmentExpression(obj)...)
	case "AwaitExpression":
		res = append(res, p.parseAwaitExpression(obj)...)
	case "ArrowFunctionExpression":
		res = append(res, p.parseArrowFunctionExpression(obj)...)
	case "SequenceExpression":
		res = append(res, p.parseSequenceExpression(obj)...)
	case "UpdateExpression":
		res = append(res, p.parseUpdateExpression(obj)...)
	case "ContinueStatement":
		res = append(res, p.parseContinueStatement(obj)...)
	case "BreakStatement":
		res = append(res, p.parseBreakStatement(obj)...)
	case "NewExpression":
		res = append(res, p.parseNewExpression(obj)...)
	case "ComputedMemberExpression":
		res = append(res, p.parseComputedMemberExpression(obj)...)
	case "MemberExpression":
		res = append(res, p.parseMemberExpression(obj)...)
	case "ReturnStatement":
		res = append(res, p.parseReturnStatement(obj)...)
	case "FunctionDeclaration":
		res = append(res, p.parseFunctionDeclaration(obj)...)
	case "ClassDeclaration":
		res = append(res, p.parseClassDeclaration(obj)...)
	case "SwitchStatement":
		res = append(res, p.parseSwitchStatement(obj)...)
	case "SwitchCase":
		res = append(res, p.parseSwitchCase(obj)...)
	case "WhileStatement":
		res = append(res, p.parseWhileStatement(obj)...)
	case "DoWhileStatement":
		res = append(res, p.parseDoWhileStatement(obj)...)
	case "ForStatement":
		res = append(res, p.parseForStatement(obj)...)
	case "ForInStatement":
		res = append(res, p.parseForInStatement(obj)...)
	case "IfStatement":
		res = append(res, p.parseIfStatement(obj)...)
	case "TryStatement":
		res = append(res, p.parseTryStatement(obj)...)
	case "MethodDefinition":
		res = append(res, p.parseMethodDefinition(obj)...)
	}
	return res
}
