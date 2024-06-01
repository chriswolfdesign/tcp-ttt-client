package client

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"

	"github.com/chriswolfdesign/tcp-ttt-common/enums"
	"github.com/chriswolfdesign/tcp-ttt-common/strings"
	"github.com/chriswolfdesign/tcp-ttt-common/tcp_payloads"
)

type Message struct {
	Name string
	Text string
}

type Client struct {
	Name       string
	Host       string
	Player     string
	ServerConn net.Conn
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

	c.ServerConn = conn

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
	if c.Player == enums.PLAYER_ONE {
		fmt.Println("You are registered, you will be playing as player 1")
	} else {
		fmt.Println("You are registered, you will be playing as player 2")
	}
}

func (c *Client) WaitForGameStart() {
	gameStarted := false

	for !gameStarted {
		gameStartedBuf := make([]byte, 1024)
		_, err := c.ServerConn.Read(gameStartedBuf)
		if err != nil {
			continue
		}

		tmp := bytes.NewBuffer(gameStartedBuf)

		gameStartedMessage := &tcp_payloads.GameStartingMessage{}
		dec := gob.NewDecoder(tmp)

		if err = dec.Decode(gameStartedMessage); err != nil {
			continue
		}

		if gameStartedMessage.PayloadType == strings.TYPE_GAME_STARTING_MESSAGE {
			gameStarted = true
		}

		fmt.Println("Game has begun")

		gameStartedMessage.Game.Board.PrintBoard()
	}
}

func (c *Client) MakeMove() {
	illegalMove := true

	for illegalMove {
		fmt.Print("Row: ")
		var rowString string
		_, err := fmt.Scanln(&rowString)
		if err != nil {
			fmt.Println(err)
			continue
		}

		row, err := strconv.Atoi(rowString)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Print("Col: ")
		var colString string
		_, err = fmt.Scanln(&colString)
		if err != nil {
			fmt.Println(err)
			continue
		}

		col, err := strconv.Atoi(colString)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var makeMoveBuffer bytes.Buffer
		enc := gob.NewEncoder(&makeMoveBuffer)

		makeMoveRequest := tcp_payloads.MakeMoveMessage{
			Row: row,
			Col: col,
			PayloadType: strings.TYPE_MAKE_MOVE_MESSAGE,
		}

		if err = enc.Encode(makeMoveRequest); err != nil {
			fmt.Println(err)
			continue
		}

		_, err = c.ServerConn.Write(makeMoveBuffer.Bytes())
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Sent move to server")

		responseBuf := make([]byte, 1024)
		_, err = c.ServerConn.Read(responseBuf)
		if err != nil {
			fmt.Println(err)
			continue
		}

		tmp := bytes.NewBuffer(responseBuf)

		response := &tcp_payloads.IllegalMoveMessage{}
		dec := gob.NewDecoder(tmp)

		if err = dec.Decode(response); err != nil {
			fmt.Println(err)
			continue
		}

		if response.PayloadType == strings.TYPE_ILLEGAL_MOVE_MESSAGE {
			fmt.Println(response.ErrorMessage)
			continue
		} else if response.PayloadType == strings.TYPE_ACCEPTED_MOVE_MESSAGE {
			illegalMove = false
		} else {
			fmt.Println("Unknown payload type:", response.PayloadType)
			continue
		}
	}
}
