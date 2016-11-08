package chatr

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// Client This type is used to connect to chatr.Server instances using a
// TCP socket. It also provides channels to be able to send data and receive
// data from the server. DO NOT CLOSE THESE CHANNELS BECAUSE THE SERVER DOES
// IT FOR US.
type Client struct {
	host     string
	port     int
	reader   *bufio.Reader
	writer   *bufio.Writer
	Incoming chan string
	Outgoing chan string
}

// read This function should be used with a go routine to have a client listen
// for incoming data from the server and then send that data to the Incoming
// channel.
func (client *Client) read() {
	// Continuously listen for data coming from network channel
	for {
		buf := make([]byte, 1024)

		n, _ := client.reader.Read(buf)
		if n == 0 {
			close(client.Incoming)
			break
		}
		message := string(buf)
		if len(strings.TrimSpace(message)) > 0 {
			client.Incoming <- message
		}
	}
}

// write This function should be used with a go routine to listen to the
// outgoing channel and send data to the server when something appears
// in the channel.
func (client *Client) write() {
	for data := range client.Outgoing {
		client.writer.WriteString(data)
		client.writer.Flush()
	}
}

// NewClient This function creates a new client with the specified host and
// port to connect to.
func NewClient(host string, port int) *Client {
	result := new(Client)
	result.host = host
	result.port = port
	result.Incoming = make(chan string)
	result.Outgoing = make(chan string)
	return result
}

// Start This function connects to the specified server and then Continuously
// waits for the server to send information. This function does not block
// because it uses go routines.
func (client *Client) Start() {
	fmt.Printf("[Client] Connecting to %v:%v\n", client.host, client.port)
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", client.host, client.port))

	if err != nil {
		log.Fatal("[Error] Failed to connect to server: ", err.Error())
	}
	client.reader = bufio.NewReader(conn)
	client.writer = bufio.NewWriter(conn)

	go client.read()
	go client.write()
}

// Send Sends a message to the connected server
func (client *Client) Send(message string) {
	client.Outgoing <- message
}
