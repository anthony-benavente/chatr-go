package chatr

import (
	"fmt"
	"log"
	"net"
	"sync"
)

// Server The EchoServer is a TCP server that listens for connections and
// when they occur, the server responds to the client with whatever the client
// sends it.
type Server struct {
	mu        sync.Mutex
	host      string
	port      int
	connected map[string]*ServerUser
	listener  net.Listener
}

// NewChatrServer This function creates a new server struct with the specified
// port. NOTE: This does not start the server until the Start method is
// invoked.
func NewChatrServer(host string, port int) (result *Server) {
	result = new(Server)
	result.host = host
	result.port = port
	result.connected = make(map[string]*ServerUser)
	return
}

// handleConnection This function handles a connection. This should be called
// using a go routine.
func (server *Server) handleConnection(conn net.Conn) {
	fmt.Printf("[Server] Got connected from %v\n", conn.RemoteAddr().String())

	// Add connection to list of connected
	newUser := NewServerUser(conn)
	newUser.OnDisconnected = func() {
		server.Broadcast(fmt.Sprintf("User %q disconnected.", newUser.conn.RemoteAddr()))
		close(newUser.incoming)
		close(newUser.outgoing)
		delete(server.connected, newUser.conn.RemoteAddr().String())
	}
	server.mu.Lock()
	server.connected[conn.RemoteAddr().String()] = newUser
	server.mu.Unlock()
	go newUser.Start()

	go func() {
		for data := range newUser.incoming {
			server.Broadcast(fmt.Sprintf("[%q] %v", newUser.conn.RemoteAddr(), data))
		}
	}()
}

// Broadcast This function sends a message to all connected clients
func (server *Server) Broadcast(message string) {
	for k := range server.connected {
		server.connected[k].outgoing <- message
	}
}

// awaitConnections This function waits for connections from the server's
// socket listener and when they connect, are sent to the handleConnection
// method using a go routine.
func (server *Server) awaitConnections() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			fmt.Printf("[Error] Server: Failed to accept connection...\n")
		} else {
			go server.handleConnection(conn)
		}
	}
}

// Start This function starts up the server at the port specified when the
// server was initialized. If the server fails to open, the program will shut
// down from a log.Fatal.
func (server *Server) Start() {
	addr := fmt.Sprintf("%v:%v", server.host, server.port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Failed to open socket at :8080")
	}
	server.listener = ln

	fmt.Println("Listening for connections on ", addr)

	server.awaitConnections()
}
