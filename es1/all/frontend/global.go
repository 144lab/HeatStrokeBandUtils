package frontend

import (
	"syscall/js"
)

var (
	window    = js.Global()
	document  = js.Global().Get("document")
	location  = js.Global().Get("location")
	navigator = js.Global().Get("navigator")
	console   = js.Global().Get("console")
)
