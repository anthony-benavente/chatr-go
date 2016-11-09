package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"flag"

	"os/user"

	chatr "github.com/anthony-benavente/chatr-go"
)

func usageQuit() {
	log.Fatal("Usage: chatr <username> -h host -p port")
}

func main() {
	if len(os.Args) < 3 {
		usageQuit()
	}

	currentUser, _ := user.Current()

	host := flag.String("h", "localhost", "-h <hostname>")
	port := flag.Int("p", 8080, "-p <port>")
	user := flag.String("u", currentUser.Username, "-u <username>")
	flag.Parse()

	stdin := bufio.NewReader(os.Stdin)
	client := chatr.NewClient(*host, *port)
	client.Start(*user)
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
