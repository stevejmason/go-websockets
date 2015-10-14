package server

import (

)

/**
 * Interface that handles client messages
 */
type ClientHandler interface {
	OnMessageReceived(client Client, message string)
}

