package server

import (
	"log"
	"net/http"
	"golang.org/x/net/websocket"
)

// - Public interface
type Server interface {
	Clients() map[uint]Client
	Listen()
	Del(c Client)
	Error(err error)
	ClientHandler() ClientHandler
	SetClientHandler(clientHandler ClientHandler)
}

// - Private struct
type server struct {
	context string
	clients map[uint]Client

	// - Client handler
	clientHandler ClientHandler

	// - Channels
	newClientCh     chan Client
	delClientCh     chan Client
	errorCh         chan error
}

// - Initializer
func NewServer(context string) Server {
	svr := &server{}

	// - Vars
	svr.context = context
	svr.clients = make(map[uint]Client)

	// - Channels
	svr.newClientCh = make(chan Client)
	svr.delClientCh = make(chan Client)
	svr.errorCh     = make(chan error)

	return svr
}

/**
 * Server constructor
 */
func (server *server) Clients() map[uint]Client {
	return server.clients
}

func (server *server) ClientHandler() ClientHandler {
	return server.clientHandler
}

func (server *server) SetClientHandler(clientHandler ClientHandler) {
	server.clientHandler = clientHandler
}

/**
 * Start server to listen & serve
 */
func (server *server) Listen() {
	log.Println("Server Listen() started...")

	// - Start listener
	http.Handle(server.context, websocket.Handler(server.onConnectionOpened))

	// - Main loop
	for {
		select {

		// - Handle new client connection
		case c := <- server.newClientCh:
			log.Println("Added new client")
			server.clients[c.Id()] = c
			server.logClients()

		case c:= <- server.delClientCh:
			log.Println("Client disconnected")
			server.stopClient(c)
			server.logClients()

		// - Handle errors
		case err := <- server.errorCh:
			log.Println("Error:", err)
		}
	}
}

/** 
 * Close connection
 */
func (server *server) Close(ws * websocket.Conn) {
	err := ws.Close()

	if(err != nil) {
		server.Error(err)
	}
}

/** 
 * Delete client
 */
func (server *server) Del(client Client) {
	log.Println("Requesting client disconnect")
	server.delClientCh <- client
}

/**
 * On error
 */
func (server *server) Error(err error) {
	server.errorCh <- err
}

/**
 * Event -> On connection opened
 */
func (server *server) onConnectionOpened(conn *websocket.Conn) {
	defer server.Close(conn)

	log.Printf("Connection opened: %v", conn.Request().RemoteAddr)
	client := NewClient(conn, server)
	server.startClient(client)
}

/**
 * Event -> New client (Client Handler Thread)
 */
func (server *server) startClient(client Client) {
	server.newClientCh <- client
	client.Listen()
}

/**
 * Stop & remove client
 */
func (server *server) stopClient(client Client) {
	delete(server.clients, client.Id())
}

/**
 * LOG -> Number of connected clients
 */
func (server *server) logClients() {
	log.Println("Now", len(server.Clients()), "clients connected.")
}


