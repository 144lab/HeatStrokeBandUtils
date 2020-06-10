package main

import (
	"log"

	"github.com/gopherjs/vecty"

	"hrm-logger/frontend"
)

func main() {
	log.SetFlags(log.Lshortfile)
	go func() {
		frontend.Setup()
		top := frontend.NewTopView()
		vecty.RenderBody(vecty.Component(top))
	}()
	select {}
}
