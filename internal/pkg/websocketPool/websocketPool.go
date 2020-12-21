package websocketpool

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketPool struct {
	pool map[string]*websocket.Conn
	mux  *sync.RWMutex
}

func NewWebsocketPool() WebsocketPool {
	return WebsocketPool{
		pool: make(map[string]*websocket.Conn),
		mux:  &sync.RWMutex{},
	}
}

func (w WebsocketPool) Add(userID string, conn *websocket.Conn) {
	w.mux.RLock()
	w.pool[userID] = conn
	w.mux.RUnlock()
}

func (w WebsocketPool) Delete(userID string) {
	w.mux.RLock()
	delete(w.pool, userID)
	w.mux.RUnlock()
}

func (w WebsocketPool) Get(userID string) (*websocket.Conn, bool) {
	conn, ok := w.pool[userID]
	return conn, ok
}
