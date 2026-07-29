package hub

import (
	"encoding/json"
	"log"
	"time"

	"github.com/XsedoX/RoomPlay/domain/room/room_id"
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/XsedoX/RoomPlay/infrastructure/client_message/client_message_publisher"
	"github.com/XsedoX/RoomPlay/infrastructure/hubs/connection_id"
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

	roomId room_id.RoomId

	connectionId connection_id.ConnectionId

	// from api
	receivedMessages chan any

	mainHub IHub

	clientMessagePublisher client_message_publisher.IClientMessagePublisher
}

func (c Client) UserId() user_id.UserId {
	return c.userId
}

func (c Client) RoomId() room_id.RoomId {
	return c.roomId
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.mainHub.UnregisterClient(&ClientRoomRequest{
			Client: c,
			RoomId: c.roomId,
		})
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
			messageBytes, bytesErr := json.Marshal(message)
			if bytesErr != nil {
				log.Printf("Failed to marshal message: %s", bytesErr)
				continue
			}
			err := c.conn.WriteMessage(websocket.TextMessage, messageBytes)
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

func (c *Client) readPump() {
	defer func() {
		c.mainHub.UnregisterClient(&ClientRoomRequest{
			Client: c,
			RoomId: c.roomId,
		})
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
		c.clientMessagePublisher.Publish(message, c.userId, c.connectionId)
	}
}

func NewClient(
	conn *websocket.Conn,
	userId user_id.UserId,
	roomId room_id.RoomId,
	clientMessagePublisher client_message_publisher.IClientMessagePublisher,
	mainHub IHub,
) *Client {
	return &Client{
		conn:                   conn,
		receivedMessages:       make(chan any, 100),
		userId:                 userId,
		roomId:                 roomId,
		clientMessagePublisher: clientMessagePublisher,
		connectionId:           connection_id.New(),
		mainHub:                mainHub,
	}
}

func (c Client) ConnectionId() connection_id.ConnectionId {
	return c.connectionId
}

func (c *Client) sendMessage(message any) {
	c.receivedMessages <- message
}

func (c *Client) shutdown() {
	close(c.receivedMessages)
	_ = c.conn.Close()
}
