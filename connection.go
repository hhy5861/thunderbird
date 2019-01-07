package thunderbird

import (
	"encoding/json"
	"github.com/hhy5861/logger"
	"log"
	"sync"
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

type Connection struct {
	tb            *Thunderbird
	ws            *websocket.Conn
	subscriptions map[string]map[string]bool
	subMutex      sync.RWMutex
	send          chan Event
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Connection) readPump() {
	defer func() {
		c.tb.disconnected(c)
		_ = c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	err := c.ws.SetReadDeadline(time.Now().Add(time.Duration(pongWait)))
	if err != nil {
		logger.Error(err)
	}

	c.ws.SetPongHandler(func(d string) error {
		err := c.ws.SetReadDeadline(time.Now().Add(time.Duration(pongWait)))
		return err
	})

	for {
		var event Event
		err := c.ws.ReadJSON(&event)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		if event.Event == "" {
			event.Event = "default"
		}

		c.tb.connected(c)
		switch event.Type {
		case "subscribe":
			c.Subscribed(event)
		case "unsubscribe":
			c.Unsubscribe(event)
			c.tb.disconnected(c)
		case "message":
			for _, ch := range c.tb.Channels(event.Channel, event.Event) {
				ch.Received(event)
			}
		default:
			log.Printf("unknown event command %s", event.Type)
		}
	}
}

func (c *Connection) Subscribed(event Event) {
	c.subMutex.Lock()
	c.subscriptions[event.Channel] = make(map[string]bool)
	c.subscriptions[event.Channel][event.Event] = true
	c.subMutex.Unlock()
}

func (c *Connection) isSubscribedTo(event Event) bool {
	c.subMutex.Lock()
	r := c.subscriptions[event.Channel][event.Event]
	c.subMutex.Unlock()

	return r
}

func (c *Connection) Unsubscribe(event Event) {
	c.subMutex.Lock()
	subscribes, ok := c.subscriptions[event.Channel]
	if ok {
		if len(subscribes) > 2 {
			delete(c.subscriptions[event.Channel], event.Event)
		} else {
			delete(c.subscriptions, event.Channel)
		}
	}

	c.subMutex.Unlock()
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.ws.Close()
	}()

	for {
		select {
		case event, ok := <-c.send:
			if !ok {
				_ = c.write(websocket.CloseMessage, []byte{})
				return
			}

			b, err := json.Marshal(event)
			if err != nil {
				return
			}

			if err := c.write(websocket.TextMessage, b); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	err := c.ws.SetWriteDeadline(time.Now().Add(time.Duration(writeWait)))
	if err != nil {
		return err
	}

	return c.ws.WriteMessage(mt, payload)
}
