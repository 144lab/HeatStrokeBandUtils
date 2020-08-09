package main

import (
	"log"
	"time"

	"mtband-logger/views"

	"github.com/nobonobo/spago"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	log.Println("unix:", time.Now().Unix())
	router := spago.NewRouter()
	router.Handle("/", func(key string) {
		log.Println(spago.GetURL())
		spago.SetTitle("Top")
		spago.RenderBody(views.NewTop())
	})
	if err := router.Start(); err != nil {
		println(err)
		spago.RenderBody(spago.NotFoundPage())
	}
	select {}
}
