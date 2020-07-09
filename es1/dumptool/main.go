package main

import (
	"log"
	"syscall/js"
	"time"

	"github.com/nobonobo/wecty"

	"dumptool/frontend"
	"dumptool/store"
)

func main() {
	posix := time.Now().Unix()
	log.Println("posix:", posix)
	wecty.AddMeta("viewport", "width=device-width,initial-scale=1")
	wecty.AddStylesheet("./assets/spectre.min.css")
	wecty.AddStylesheet("./assets/spectre-icons.min.css")
	wecty.AddStylesheet("./assets/spectre-exp.min.css")
	wecty.AddStylesheet("./assets/app.css")
	wecty.LoadScript("./assets/recorder.js")
	store.Recorder = js.Global().Get("HrmRecorder").New()
	wecty.RenderBody(&frontend.Top{})
	select {}
}
