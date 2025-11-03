package ws

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"server/internal/domain"
	"server/internal/server/middleware"
	"server/internal/service"
	"server/util"
)

var websocketUpgrader = &websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

type Manager struct {
	sync.RWMutex
	Clients             ClientList
	Handlers            map[string]EventHandler
	NotificationService *service.NotificationService
}

func NewManager(
	notificationService *service.NotificationService,
) *Manager {
	m := &Manager{
		Clients:             make(ClientList),
		Handlers:            make(map[string]EventHandler),
		NotificationService: notificationService,
	}

	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.Handlers[EventFavoritedListingNotification] = handleFavoritedListing
	m.Handlers[EventPriceDropNotification] = handlePriceDrop
	m.Handlers[EventStatusChangeNotification] = handleStatusChange
}

func (m *Manager) StartWSConn(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(middleware.UserContextKey).(*domain.ContextSessionData)

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		util.RespondWithError(
			w,
			http.StatusInternalServerError,
			"Error upgrading websocket connection",
		)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := NewWSClient(conn, userCtx.UserID, userCtx.Role, m)
	m.AddClient(client)

	go client.ReadMessages(ctx, cancel)
	go client.WriteMessages(ctx, cancel)
}

func (m *Manager) RouteEvent(event Event, client *WSClient) error {
	if handler, ok := m.Handlers[event.Type]; ok {
		if err := handler(event, client); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("Could not find event type")
	}
}

func (m *Manager) AddClient(client *WSClient) {
	m.Lock()
	defer m.Unlock()

	m.Clients[client] = true
}

func (m *Manager) RemoveClient(client *WSClient) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.Clients[client]; ok {
		client.Connection.Close()
		delete(m.Clients, client)
	}
}
