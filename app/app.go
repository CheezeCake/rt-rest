package app

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CheezeCake/rt-rest/config"
	"github.com/CheezeCake/rt-rest/rtorrent"
	"github.com/CheezeCake/rt-rest/rtorrent/rpc"
	"github.com/CheezeCake/rt-rest/web"
)

var (
	client rtorrent.Client
)

func Init(cfg config.Cfg) {
	var network rtorrent.Network
	if cfg.Rtorrent.Unix {
		network = rtorrent.Unix
	} else {
		network = rtorrent.Tcp
	}

	client = rtorrent.Client{Network: network, Address: cfg.Rtorrent.Address}

	setupRoutes()
}

func Test(w http.ResponseWriter, r *http.Request) error {
	err := client.Test()
	var errString string
	if err != nil {
		errString = err.Error()
	}

	json.NewEncoder(w).Encode(&struct {
		Ok    bool   `json:"ok"`
		Error string `json:"error,omitempty"`
	}{err == nil, errString})
	return nil
}

func TorrentAdd(w http.ResponseWriter, r *http.Request) error {
	var data struct {
		Uri string `json:"uri"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}
	if len(data.Uri) == 0 {
		return web.Error{errors.New("missing `uri`"), http.StatusBadRequest}
	}

	return singleCall("load.normal", []string{"", data.Uri})
}

type torrent struct {
	Hash string // d.hash
	Name string // d.name

	IsOpen   bool // d.is_open
	IsActive bool // d.is_active

	BytesDone int64 // d.bytes_done
	SizeBytes int64 // d.size_bytes

	DownRate int64
	UpRate   int64
}

func TorrentList(w http.ResponseWriter, r *http.Request) error {
	var torrents []torrent
	return multiCall("d.multicall2",
		[]string{
			"",
			"main",
			"d.hash=",
			"d.name=",
			"d.is_open=",
			"d.is_active=",
			"d.bytes_done=",
			"d.size_bytes=",
			"d.down.rate=",
			"d.up.rate=",
		},
		&torrents, w)
}

func TorrentStart(hash string, w http.ResponseWriter) error {
	return singleCall("d.start", []string{hash})
}

func TorrentStop(hash string, w http.ResponseWriter) error {
	return singleCall("d.stop", []string{hash})
}

func TorrentPause(hash string, w http.ResponseWriter) error {
	return singleCall("d.pause", []string{hash})
}

type file struct {
	Path      string
	Priority  int64
	SizeBytes int64
	Id        int
}

func TorrentFiles(hash string, w http.ResponseWriter) error {
	var files []file
	return multiCall("f.multicall",
		[]string{
			hash,
			"",
			"f.path=",
			"f.priority=",
			"f.size_bytes=",
		},
		&files,
		w)
}

type peer struct {
	Address   string
	Client    string
	DownRate  int64
	DownTotal int64
	UpRate    int64
	UpTotal   int64
	Encrypted bool
}

func TorrentPeers(hash string, w http.ResponseWriter) error {
	var peers []peer
	return multiCall("p.multicall",
		[]string{
			hash,
			"",
			"p.address=",
			"p.client_version=",
			"p.down_rate=",
			"p.down_total=",
			"p.up_rate=",
			"p.up_total=",
			"p.is_encrypted=",
		},
		&peers,
		w)
}

func singleCall(command string, args []string) error {
	res, err := client.Do(command, args)
	if err != nil {
		return err
	}
	return rpc.DecodeResponseForFault(res)
}

func multiCall(command string, args []string, result interface{}, w http.ResponseWriter) error {
	rpcRes, err := client.Do(command, args)
	if err != nil {
		return err
	}

	err = rpc.DecodeResponse(rpcRes, result)
	if err != nil {
		return err
	}
	return web.WriteJsonResponse(w, result)
}
