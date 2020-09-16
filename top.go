package main

import (
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> ce48c93... improve
	"fmt"
	"go/format"
	"strings"
=======
>>>>>>> cdc3d9a... es5 support completed
=======
	"strings"
>>>>>>> e15b2fb... improve
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

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> ce48c93... improve
func (c *Top) parse(s string) (res js.Value, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s", e)
			res = js.Null()
		}
	}()
	return esprima.Call("parseScript", s), nil
}

<<<<<<< HEAD
=======
>>>>>>> e15b2fb... improve
=======
>>>>>>> ce48c93... improve
// OnSubmit ...
func (c *Top) OnSubmit(ev js.Value) {
	ev.Call("preventDefault")
	params := jsutil.Form2Go(ev.Get("target"))
<<<<<<< HEAD
<<<<<<< HEAD
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
=======
=======
	c.JsCode = (params["js"]).(string)
<<<<<<< HEAD
>>>>>>> e15b2fb... improve
	tree := esprima.Call("parseScript", params["js"])
=======
	tree, err := c.parse(c.JsCode)
	if err != nil {
		js.Global().Call("alert", err.Error())
		return
	}
>>>>>>> ce48c93... improve
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
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> cdc3d9a... es5 support completed
=======
	c.GoCode = strings.Join(res, "\n")
=======
	c.GoCode = string(b)
>>>>>>> ce48c93... improve
	spago.Rerender(c)
>>>>>>> e15b2fb... improve
}
