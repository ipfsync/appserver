package appserver

import "github.com/gorilla/websocket"

type MessageError struct {
	Code    int
	Message string
}

type MessageCmd struct {
	Id   string
	Cmd  string
	Data map[string]string
}

type MessageReply struct {
	Id    string
	Ok    bool
	Data  map[string]string
	Error MessageError
}

type wsClient struct {
	srv  *AppServer
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *wsClient) readPump() {
	defer func() {
		c.srv.unregisterWsClient(c)
		c.conn.Close()
	}()
}
