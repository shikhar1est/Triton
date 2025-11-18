package client

import (
	"log"
	"net/http"
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
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan *Message
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		// On exit, unregister the client and close the connection
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// Set a 'pong' handler to reset the read deadline
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait));
		return nil
	})

	// This is the loop that reads messages from the client
	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Client error: %v", err)
			} else {
				log.Println("Client disconnected")
			}
			break // Exit loop on error
		}
		
		// Send the received message to the hub's broadcast channel
		c.hub.broadcast <- &msg
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	// A Ticker to send 'ping' messages
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			// Set a write deadline
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			// Write the JSON message to the connection
			err := c.conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error writing JSON to client: %v", err)
				return
			}

		case <-ticker.C:
			// This is the 'ping' message to keep the connection alive
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Error sending ping:", err)
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create a new client
	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan *Message, 256), // 256 is the buffer size
	}
	
	// Register this new client with the hub
	client.hub.register <- client

	// Start the client's goroutines
	// These run concurrently
	go client.writePump()
	go client.readPump()

	log.Println("New client connected")
}