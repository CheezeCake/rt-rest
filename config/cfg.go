package config

import (
	"encoding/json"
	"log"
	"os"
)

type rtorrentCfg struct {
	Unix    bool   `json:"unix"`
	Address string `json:"address"`
}

type Cfg struct {
	Rtorrent          rtorrentCfg `json:"rtorrent"`
	ListenningAddress string      `json:"listenningAddress"`
	ListenningPort    int         `json:"listenningPort"`
	Secret            string      `json:"secret"`
}

func Load(filename string) (Cfg, error) {
	// default config
	config := Cfg{
		Rtorrent: rtorrentCfg{
			Unix:    false,
			Address: "localhost:5000",
		},
		ListenningAddress: "127.0.0.1",
		ListenningPort:    8080,
	}

	f, err := os.Open(filename)
	if err != nil {
		return config, err
	}

	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		log.Println("Error reading config:", err)
	}
	return config, nil
}
