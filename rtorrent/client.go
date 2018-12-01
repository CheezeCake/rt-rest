package rtorrent

import (
	"bytes"
	"io"
	"net"

	"github.com/CheezeCake/rt-rest/rtorrent/rpc"
	"github.com/CheezeCake/rt-rest/scgi"
)

type Network string

const (
	Tcp  = "tcp"
	Unix = "unix"
)

type Client struct {
	Network Network
	Address string
}

func (c Client) Test() error {
	conn, err := net.Dial(string(c.Network), c.Address)
	if err == nil {
		conn.Close()
	}
	return err
}

func (c Client) Do(method string, params []string) ([]byte, error) {
	conn, err := net.Dial(string(c.Network), c.Address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	xmlRequest := rpc.EncodeRequest(method, params)
	_, err = io.Copy(conn, bytes.NewReader(scgi.EncodeRequest(xmlRequest)))
	if err != nil {
		return nil, err
	}

	var res bytes.Buffer
	_, err = io.Copy(&res, conn)
	return res.Bytes(), err
}
