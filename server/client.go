package server

import (
	"io"
	"log"
	"golang.org/x/net/websocket"
	"fmt"
)

// - Public interface
type Client interface {
	Id() uint
	Listen()
	SendMessage(message string)
}

// - Private struct
type client struct {
	id uint
	ws *websocket.Conn
	svr Server
	done bool

	// - Channels
	doneCh chan bool
	msgRecvCh chan string
	msgSendCh chan string
}

// - Variables
var maxId uint = 0

/** 
 * Client constructor
 */
func NewClient(ws *websocket.Conn, svr Server) Client {
	c := &client{}

	// - Client vars
	c.id  = maxId; maxId++
	c.ws  = ws
	c.svr = svr

	// - Channels
	c.doneCh 	= make(chan bool)
	c.msgRecvCh = make(chan string)
	c.msgSendCh = make(chan string)

	return c
}

/** 
 * GET-> id
 */
func (c *client) Id() uint {
	return c.id
}

/**
 * Queue message for sending
 */
func (c *client) SendMessage(message string) {
	c.msgSendCh <- message
}

/**
 * Handle a message
 */
func (c *client) handleMessage(msg string) {
	c.svr.ClientHandler().OnMessageReceived(c, msg)
}

/**
 * Send a message (internal)
 */
func (c *client) sendMessage(msg string) {
	log.Printf("Sending message (%s) to client (%s)\n", msg, c)
	websocket.Message.Send(c.ws, msg)
}

/** 
 * Listen main method
 */
func (c *client) Listen() {
	go c.listenWrite()
	go c.listenCallback()
	c.listenRead()
}

/**
 * Listen read 
 */
func (c *client) listenRead() {
	log.Println("Listening read from client")

	for {
		if(c.done) {
			return
		}

		select {

		// read data from websocket connection
		default:
			var msg string
		//			err := websocket.JSON.Receive(c.ws, &msg)
			err := websocket.Message.Receive(c.ws, &msg)
			if err == io.EOF {
				c.disconnect()
			} else if err != nil {
				c.svr.Error(err)
				c.disconnect()
			} else {
				// - Handle message
				log.Println("Received message:", msg)
				c.msgRecvCh <- msg
			}
		}
	}
}

/**
 * Listen write 
 */
func (c *client) listenWrite() {
	for {
		if(c.done) {
			return
		}

		select {
			// - Send message
			case msg := <- c.msgSendCh:
				c.sendMessage(msg)
		}
	}
}

/**
 * Callback listener
 */
func (c *client) listenCallback() {
	for {
		if(c.done) {
			return
		}

		select {

			// - Post message received callbacks
			case msg := <- c.msgRecvCh:
				c.handleMessage(msg)

			// -
		}
	}
}

/**
 * Disconnect client
 */
func (c *client) disconnect() {
	// - Skip if disconnected already
	log.Println("Client disconnecting", c.id)

	// - Close
	err := c.ws.Close()
	if(err != nil) {
		c.svr.Error(err)
	}

	// - Remove from server
	c.svr.Del(c)

	// - Mark done for any alive read/writers
	c.doneCh <- true
}

/**
 * Dump information as string
 */
func (c *client) String() string {
	return fmt.Sprintf("Client -> ID: %d , remote: %v", c.id, c.ws.Request().RemoteAddr)
}
