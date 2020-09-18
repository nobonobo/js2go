package main

import (
	"fmt"
	"go/format"
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
	GoCode string
}

func (c *Top) parse(s string) (res js.Value, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s", e)
			res = js.Null()
		}
	}()
	return esprima.Call("parseScript", s), nil
}

// OnSubmit ...
func (c *Top) OnSubmit(ev js.Value) {
	ev.Call("preventDefault")
	params := jsutil.Form2Go(ev.Get("target"))
	c.JsCode = (params["js"]).(string)
	tree, err := c.parse(c.JsCode)
	if err != nil {
		js.Global().Call("alert", err.Error())
		return
	}
	console.Call("log", tree)
	parser := &Parser{}
	res, err := parser.ParseProgram(tree)
	if err != nil {
		js.Global().Call("alert", err.Error())
		return
	}
	b, err := format.Source([]byte(strings.Join(res, "\n")))
	if err != nil {
		js.Global().Call("alert", err.Error())
		return
	}
	c.GoCode = string(b)
	spago.Rerender(c)
}
