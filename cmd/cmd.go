package main

import (
	"fmt"
	"tcp-ttt-client/client"
)

func main() {
	playerName := getPlayerName()

	cl := client.GenerateClient(playerName, "localhost:8080")
	cl.RegisterPlayer()
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
