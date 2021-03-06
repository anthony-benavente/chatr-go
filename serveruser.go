package chatr

import (
	"bufio"
	"net"
)

// ServerUser This struct represents users that are connected. This gives
// us an interface to talk to connected clients.
type ServerUser struct {
	username       string
	conn           net.Conn
	incoming       chan string
	outgoing       chan string
	OnDisconnected func()
}

// NewChatrServerUser Creates a new user in the chatr server with the given
// network connection
func NewServerUser(username string, conn net.Conn) *ServerUser {
	result := new(ServerUser)
	result.username = username
	result.conn = conn
	result.incoming = make(chan string)
	result.outgoing = make(chan string)
	result.OnDisconnected = func() {}
	return result
}

// Start This function starts a go routine waiting to send data out to the
// connected client
func (user *ServerUser) Start() {
	reader := bufio.NewReader(user.conn)
	writer := bufio.NewWriter(user.conn)
	go func() {
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				user.OnDisconnected()
				break
			} else {
				user.incoming <- msg
			}
		}
	}()

	for data := range user.outgoing {
		writer.WriteString(data)
		writer.Flush()
	}
}
