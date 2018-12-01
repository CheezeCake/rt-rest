package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/CheezeCake/rt-rest/app"
	"github.com/CheezeCake/rt-rest/config"
	"github.com/CheezeCake/rt-rest/data"
	"github.com/CheezeCake/rt-rest/web"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	app.Init(cfg)
	data.Init(cfg)
	web.Init(cfg)

	http.Handle("/", http.FileServer(http.Dir("/home/manu/vue_test")))
	log.Fatal(http.ListenAndServe(cfg.ListenningAddress+":"+strconv.Itoa(cfg.ListenningPort), nil))
}
