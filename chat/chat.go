/**
 * Server handler implementation that broadcasts messages to all other clients
 */
package chat

import (
	"log"
	"github.com/stevejmason/go-websockets/server"
)

type ChatServer interface {
	server.ClientHandler
}

type chatServer struct {
	svr server.Server
}

/**
 * Constructor
 */
func NewChatServer(svr server.Server) ChatServer {
	return &chatServer{svr}
}

func (csvr *chatServer) OnMessageReceived(client server.Client, message string) {
	log.Println("Chat server got message:", message, "from client:", client)
	csvr.broadcastMessage(message)
}

func (csvr *chatServer) broadcastMessage(message string, excludeClients ... server.Client) {
	for _, client := range csvr.svr.Clients() {
		client.SendMessage(message)
	}
}
