package main

import (
	"log"
	"net/http"
	"github.com/stevejmason/go-websockets/server"
	"github.com/stevejmason/go-websockets/chat"
)

const WS_CONTEXT = "/server"

func main() {
	log.SetFlags(log.Lshortfile)

	// - Create server
	svr := server.NewServer(WS_CONTEXT)

	// - Chat server handler
	chatHandler := chat.NewChatServer(svr)

	// - Setup svr
	svr.SetClientHandler(chatHandler)

	// - Debug
	log.Println("Server object", svr)

	// - Listen
	go svr.Listen()

	// - Static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	// - Start http listener
	log.Fatal(http.ListenAndServe(":8080", nil))
}
