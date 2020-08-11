package main

import (
	"log"
	"time"

	"mtband-logger/views"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/router"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	log.Println("unix:", time.Now().Unix())
	r := router.New()
	r.Handle("/", func(key string) {
		log.Println(router.GetURL())
		spago.SetTitle("Top")
		spago.RenderBody(views.NewTop())
	})
	if err := r.Start(); err != nil {
		println(err)
		spago.RenderBody(router.NotFoundPage())
	}
	select {}
}
