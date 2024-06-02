package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"tcp-ttt-client/client"

	"github.com/chriswolfdesign/tcp-ttt-common/model"
	"github.com/chriswolfdesign/tcp-ttt-common/strings"
	"github.com/chriswolfdesign/tcp-ttt-common/tcp_payloads"
)

func main() {
	playerName := getPlayerName()

	cl := client.GenerateClient(playerName, "localhost:8080")
	cl.RegisterPlayer()

	fmt.Println("Waiting for game to begin")
	gameStarted := false

	var game *model.Game

	for !gameStarted {
		gameStartedBuf := make([]byte, 1024)
		_, err := cl.ServerConn.Read(gameStartedBuf)
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
			game = &gameStartedMessage.Game
			gameStarted = true
		}

		fmt.Println("Game has begun")

		gameStartedMessage.Game.Board.PrintBoard()
	}

	gameOver := false

	for !gameOver {

		if game.Winner != strings.NOT_OVER {
			fmt.Println("The game is over")
			fmt.Println("Result:", game.Winner)
			return
		}

		if game.CurrentPlayer == cl.Player {
			cl.MakeMove()
		} else {
			fmt.Println("Waiting for other player")
		}

		gameStateBuf := make([]byte, 1024)
		_, err := cl.ServerConn.Read(gameStateBuf)
		if err != nil {
			fmt.Println(err)
			continue
		}

		tmp := bytes.NewBuffer(gameStateBuf)

		gameStateResponse := &tcp_payloads.GameStateMessage{}

		dec := gob.NewDecoder(tmp)

		if err = dec.Decode(gameStateResponse); err != nil {
			fmt.Println(err)
			continue
		}

		game = &gameStateResponse.Game
		game.Board.PrintBoard()
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
