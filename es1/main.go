package main

import (
	"log"
	"time"

	"github.com/nobonobo/wecty"

	"es1/views"
)

func main() {
	posix := time.Now().Unix()
	log.Println("posix:", posix)
	wecty.AddStylesheet("./assets/spectre.min.css")
	wecty.AddStylesheet("./assets/spectre-icons.min.css")
	wecty.AddStylesheet("./assets/spectre-exp.min.css")
	wecty.AddStylesheet("./assets/app.css")
	wecty.LoadScript("./assets/recorder.js")
	wecty.RenderBody(views.NewTop())
	select {}
}
