package thunderbird

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type (
	Adapter interface {
		Subscribe(event Event) error
		Unsubscribe(event Event) error
		Broadcast(event Event, payload []byte) error
	}
)

func New() *Thunderbird {
	return &Thunderbird{
		connections:     make(map[*Connection]bool),
		channelHandlers: make(map[string]map[string][]ChannelHandler),
	}
}

type Thunderbird struct {
	channelHandlers map[string]map[string][]ChannelHandler
	chanMutex       sync.RWMutex
	connections     map[*Connection]bool
	connMutex       sync.RWMutex
	openSend        bool
}

func (tb *Thunderbird) SetOpenSend(openSend bool) {
	tb.openSend = openSend
}

func (tb *Thunderbird) Broadcast(event Event) {
	tb.connMutex.Lock()
	for conn := range tb.connections {
		if conn.isSubscribedTo(event) {
			event := Event{
				Type:    "message",
				Channel: event.Channel,
				Event:   event.Event,
				Body:    event.Body,
			}

			conn.send <- event
		}
	}

	tb.connMutex.Unlock()
}

func (tb *Thunderbird) newConnection(ws *websocket.Conn) *Connection {
	return &Connection{
		tb:            tb,
		subscriptions: make(map[string]map[string]bool),
		ws:            ws,
		send:          make(chan Event),
	}
}

func (tb *Thunderbird) HTTPHandler() http.Handler {
	return &httpHandler{
		tb: tb,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (tb *Thunderbird) HTTPHandlerWithUpgrader(upgrader websocket.Upgrader) http.Handler {
	return &httpHandler{tb: tb, upgrader: upgrader}
}

func (tb *Thunderbird) connected(c *Connection) {
	tb.connMutex.Lock()
	tb.connections[c] = true
	tb.connMutex.Unlock()
}

func (tb *Thunderbird) subscribed(e Event) {
}

func (tb *Thunderbird) disconnected(c *Connection) {
	tb.connMutex.Lock()
	if ok := tb.connections[c]; ok {
		delete(tb.connections, c)
		close(c.send)
	}

	tb.connMutex.Unlock()
}

func (tb *Thunderbird) HandleChannel(channel, event string, handler ChannelHandler) {
	tb.chanMutex.Lock()
	tb.channelHandlers[channel] = make(map[string][]ChannelHandler)
	tb.channelHandlers[channel][event] = append(tb.channelHandlers[channel][event], handler)
	tb.chanMutex.Unlock()
}

func (tb *Thunderbird) Channels(channel, event string) []ChannelHandler {
	tb.chanMutex.Lock()
	ch := tb.channelHandlers[channel][event]
	tb.chanMutex.Unlock()

	return ch
}

type ChannelHandler interface {
	Received(Event)
}
