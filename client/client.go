package client

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"

	"github.com/chriswolfdesign/tcp-ttt-common/tcp_payloads"
)

type Message struct {
	Name string
	Text string
}

type Client struct {
	Name   string
	Host   string
	Player string
}

func GenerateClient(name, host string) Client {
	return Client{
		Name: name,
		Host: host,
	}
}

func (c *Client) RegisterPlayer() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	onboardingRequest := tcp_payloads.GeneratePlayerOnboardingRequest(c.Name)

	if err = enc.Encode(onboardingRequest); err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}

	responseBuf := make([]byte, 1024)
	_, err = conn.Read(responseBuf)
	if err != nil {
		fmt.Println(err)
		return
	}

	tmp := bytes.NewBuffer(responseBuf)

	response := &tcp_payloads.PlayerOnboardingResponse{}
	dec := gob.NewDecoder(tmp)

	if err = dec.Decode(response); err != nil {
		fmt.Println(err)
		return
	}

	c.Player = response.Player
	fmt.Printf("Client: %+v\n", c)

	conn.Close()
}
