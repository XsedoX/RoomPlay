package room_hub

import (
	"log"
	"time"

	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/client_message/client_message_publisher"
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

type Client struct {
	conn *websocket.Conn

	userId user_id.UserId

	// from api
	receivedMessages chan []byte

	roomHub *RoomHub

	clientMessagePublisher client_message_publisher.IClientMessagePublisher
}

func (c *Client) UserId() user_id.UserId {
	return c.userId
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.conn.Close()
		ticker.Stop()
		c.roomHub.UnregisterClient(c)
	}()

	for {
		select {
		case message, ok := <-c.receivedMessages:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.roomHub.UnregisterClient(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.clientMessagePublisher.Publish(message, c.userId)
	}
}

func NewClient(
	conn *websocket.Conn,
	userId user_id.UserId,
	clientMessagePublisher client_message_publisher.IClientMessagePublisher,
) *Client {
	return &Client{
		conn:                   conn,
		receivedMessages:       make(chan []byte, 100),
		userId:                 userId,
		clientMessagePublisher: clientMessagePublisher,
	}
}

func (c *Client) SendMessage(message []byte) {
	c.receivedMessages <- message
}

func (c *Client) Close() {
	c.conn.Close()
	close(c.receivedMessages)
}
