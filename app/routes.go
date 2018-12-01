package app

import (
	"net/http"

	"github.com/CheezeCake/rt-rest/web"
)

func setupRoutes() {
	// http.HandleFunc("/api/test",
	// 	web.NewChain(web.Method("GET"), web.Auth()).Finally(Test))
	http.HandleFunc("/api/test",
		web.NewChain(web.Method("GET")).Finally(Test))

	http.HandleFunc("/api/login",
		web.NewChain(web.Method("POST")).Finally(web.Login))
	http.HandleFunc("/api/register",
		web.NewChain(web.Method("POST")).Finally(web.Register))

	http.HandleFunc("/api/torrent/add",
		web.NewChain(web.Method("POST"), web.Auth()).Finally(TorrentAdd))

	http.HandleFunc("/api/torrent/list",
		web.NewChain(web.Method("GET"), web.Auth()).Finally(TorrentList))

	http.HandleFunc("/api/torrent/start",
		web.NewChain(web.Method("PATCH"), web.Auth()).Finally(web.MakeTorrentHandler(TorrentStart)))
	http.HandleFunc("/api/torrent/stop",
		web.NewChain(web.Method("PATCH"), web.Auth()).Finally(web.MakeTorrentHandler(TorrentStop)))
	http.HandleFunc("/api/torrent/pause",
		web.NewChain(web.Method("PATCH"), web.Auth()).Finally(web.MakeTorrentHandler(TorrentPause)))

	http.HandleFunc("/api/torrent/files",
		// web.NewChain(web.Method("GET"), web.Auth()).Finally(web.MakeTorrentHandler(TorrentFiles)))
		web.NewChain(web.Method("GET")).Finally(web.MakeTorrentHandler(TorrentFiles)))
}
