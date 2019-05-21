package appserver

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

const (
	errCodeUnknownCmd = 1000 // Unknown command error
	errCodeInternal   = 1001 // Internal error
)

type wsClient struct {
	srv  *AppServer
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan interface{}
}

func newWsClient(srv *AppServer, conn *websocket.Conn) *wsClient {
	return &wsClient{
		srv:  srv,
		conn: conn,
		send: make(chan interface{}, 100),
	}
}

func (c *wsClient) readPump() {
	defer func() {
		c.srv.unregisterWsClient(c)
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var msg MessageCmd
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Invalid cmd error: %v", err)
			break
		}
		go c.handleCmd(&msg)
	}
}

func (c *wsClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.srv.unregisterWsClient(c)
		_ = c.conn.Close()
	}()

	for {
		select {

		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(message)
			if err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}

func (c *wsClient) sendReply(id string, data map[string]interface{}) {
	msg := &MessageReply{
		Id:   id,
		Ok:   true,
		Data: data,
	}
	c.send <- msg
}

func (c *wsClient) sendError(id string, errCode int, errMsg string, data map[string]interface{}) {
	msg := &MessageReply{
		Id:    id,
		Ok:    false,
		Data:  data,
		Error: MessageError{errCode, errMsg},
	}
	c.send <- msg
}

func (c *wsClient) handleCmd(msg *MessageCmd) {
	switch msg.Cmd {
	case "peers":
		peers, _, err := c.srv.api.Peers()
		if err != nil {
			c.sendError(msg.Id, errCodeInternal, err.Error(), nil)
			return
		}
		c.sendReply(msg.Id, map[string]interface{}{"peers": peers})
	default:
		c.sendError(msg.Id, errCodeUnknownCmd, "Unknown Command", nil)
	}
}
