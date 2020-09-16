package main

import (
	"log"
	"syscall/js"

	"github.com/nobonobo/spago"
)

var (
	esprima = js.Global().Get("esprima")
	console = js.Global().Get("console")
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	spago.RenderBody(&Top{JsCode: src})
	select {}
}
