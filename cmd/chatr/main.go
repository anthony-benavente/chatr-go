package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"strings"

	"strconv"

	chatr "github.com/anthony-benavente/chatr-go"
)

func usageQuit() {
	log.Fatal("Usage: chatr <username> <host>:<port>")
}

func main() {
	if len(os.Args) < 3 {
		usageQuit()
	}

	var host string
	var port int

	if addr := strings.Split(os.Args[2], ":"); len(addr) == 2 {
		host = addr[0]
		if pport, err := strconv.Atoi(addr[1]); err != nil {
			usageQuit()
		} else {
			port = pport
		}
	}

	stdin := bufio.NewReader(os.Stdin)
	client := chatr.NewClient(host, port)
	client.Start(os.Args[1])
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
