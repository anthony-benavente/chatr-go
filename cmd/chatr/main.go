package main

import (
	"bufio"
	"fmt"
	"os"

	chatr "github.com/anthony-benavente/chatr-go"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)

	client := chatr.NewClient("localhost", 8080)
	client.Start()

	go func() {
		for data := range client.Incoming {
			fmt.Println(data)
		}
	}()

	for {
		msg, _ := stdin.ReadString('\n')
		client.Send(msg)
	}
}
