package main

import (
	"log"

	"github.com/gopherjs/vecty"

	"hrm-logger/frontend"
)

func main() {
	log.SetFlags(log.Lshortfile)
	frontend.Setup()
	top := frontend.NewTopView()
	vecty.RenderBody(top)
	select {}
}
