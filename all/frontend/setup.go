package frontend

import (
	"strings"
	"syscall/js"

	"github.com/gopherjs/vecty"
)

var (
	window    = js.Global()
	document  = js.Global().Get("document")
	location  = js.Global().Get("location")
	navigator = js.Global().Get("navigator")
	console   = js.Global().Get("console")
)

// Setup ...
func Setup() {
	console.Call("log", navigator.Get("userAgent").String())
	if strings.Index(navigator.Get("userAgent").String(), "Chrome") < 0 {
		window.Call("alert", "This Page is Google Chrome base browser only!")
		return
	}
	meta := document.Call("createElement", "meta")
	meta.Call("setAttribute", "name", "viewport")
	meta.Call("setAttribute", "content", "width=device-width,initial-scale=1")
	document.Get("head").Call("append", meta)
	vecty.AddStylesheet("css/spectre.min.css")
	vecty.AddStylesheet("css/spectre-exp.min.css")
	vecty.AddStylesheet("css/spectre-icons.min.css")
	vecty.AddStylesheet("css/app.css")
}
