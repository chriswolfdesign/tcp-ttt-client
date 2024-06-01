package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"tcp-ttt-client/client"

	"github.com/chriswolfdesign/tcp-ttt-common/strings"
	"github.com/chriswolfdesign/tcp-ttt-common/tcp_payloads"
)

func main() {
	playerName := getPlayerName()

	cl := client.GenerateClient(playerName, "localhost:8080")
	cl.RegisterPlayer()

	fmt.Println("Waiting for game to begin")
	cl.WaitForGameStart()

	gameOver := false

	for !gameOver {
		fmt.Println("Game not over")

		playerTurnBuf := make([]byte, 1024)
		_, err := cl.ServerConn.Read(playerTurnBuf)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		tmp := bytes.NewBuffer(playerTurnBuf)

		playerTurnMessage := &tcp_payloads.PlayerTurnMessage{}
		dec := gob.NewDecoder(tmp)

		if err = dec.Decode(playerTurnMessage); err != nil {
			continue
		}

		if playerTurnMessage.PayloadType == strings.TYPE_PLAYER_TURN_MESSAGE {
			fmt.Printf("Player turn is %s!\n", playerTurnMessage.Player)
		}

		return
	}
}

func getPlayerName() string {
	var playerName string

	for playerName == "" {

		fmt.Print("Enter player name: ")
		_, err := fmt.Scanln(&playerName)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return playerName
}
