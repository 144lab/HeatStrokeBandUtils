package main

import (
	"log"

	"hrm-logger/frontend"

	"github.com/nobonobo/wecty"
)

func main() {
	log.SetFlags(log.Lshortfile)
	wecty.AddMeta("viewport", "width=device-width,initial-scale=1")
	wecty.AddStylesheet("./assets/css/spectre.min.css")
	wecty.AddStylesheet("./assets/css/spectre-icons.min.css")
	wecty.AddStylesheet("./assets/css/spectre-exp.min.css")
	wecty.AddStylesheet("./assets/css/app.css")
	wecty.LoadScript("./assets/nosleep.min.js")
	wecty.LoadScript("./assets/recorder.js")
	top := frontend.NewTopView()
	wecty.RenderBody(top)
	select {}
}
