package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

type Message struct {
	Name string
	Text string
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	message := Message {
		Name: "Chris",
		Text: "Hello",
	}

	if err = enc.Encode(message); err != nil {
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

	response := &Message{}
	dec := gob.NewDecoder(tmp)

	if err = dec.Decode(response); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Received: %+v\n", response)

	conn.Close()
}
