package main

import (
	"tcp-ttt-client/client"
)

func main() {
	cl := client.GenerateClient("Chris", "localhost:8080")
	cl.RegisterPlayer()
}
