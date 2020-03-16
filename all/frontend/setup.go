package frontend

import (
	"log"
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
	ch := make(chan bool)
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
	window.Call("addEventListener", "DOMContentLoaded",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			script := document.Call("createElement", "script")
			script.Set("src", "nosleep.min.js")
			script.Call("addEventListener", "load",
				js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					window.Set("noSleep", window.Get("NoSleep").New())
					return nil
				}),
			)
			document.Get("body").Call("appendChild", script)
			script = document.Call("createElement", "script")
			script.Set("src", "recorder.js")
			script.Call("addEventListener", "load",
				js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					close(ch)
					return nil
				}),
			)
			log.Print("contentloaded")
			document.Get("body").Call("appendChild", script)
			return nil
		}),
	)
	//document.Get("body").Call("appendChild", script)
	//log.Print("contentloaded")
	<-ch
}
