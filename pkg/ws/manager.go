package ws

import (
	"sync"
)

var (
	Manager *WSManager
	once    sync.Once
)

func InitManager() {
	once.Do(func() {
		Manager = NewManager()
		go Manager.Run()
	})
}

type WSManager struct {
	Clients    map[uint]map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mutex      sync.RWMutex
}

func NewManager() *WSManager {
	return &WSManager{
		Clients:    make(map[uint]map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (m *WSManager) Run() {
	for {
		select {
		case client := <-m.Register:
			m.mutex.Lock()
			if _, ok := m.Clients[client.ID]; !ok {
				m.Clients[client.ID] = make(map[*Client]bool)
			}
			m.Clients[client.ID][client] = true
			m.mutex.Unlock()

		case client := <-m.Unregister:
			m.mutex.Lock()
			if _, ok := m.Clients[client.ID]; ok {
				if _, ok := m.Clients[client.ID][client]; ok {
					delete(m.Clients[client.ID], client)
					close(client.Send)
				}
			}
			m.mutex.Unlock()

		case message := <-m.Broadcast:
			m.mutex.RLock()
			for _, clients := range m.Clients {
				for client := range clients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			m.mutex.RUnlock()
		}
	}
}
