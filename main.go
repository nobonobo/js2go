package main

import (
	"syscall/js"

	"github.com/nobonobo/spago"
)

var (
	esprima = js.Global().Get("esprima")
	console = js.Global().Get("console")
)

func main() {
	spago.RenderBody(&Top{JsCode: src})
	select {}
}
