package main

import (
	chatr "github.com/anthony-benavente/chatr-go"
)

func main() {
	server := chatr.NewChatrServer("0.0.0.0", 8080)
	server.Start()
}
