package bot

import (
	"runtime/debug"
	"sync"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

var _ EventManager = (*eventManagerImpl)(nil)

// NewEventManager returns a new EventManager with the EventManagerConfigOpt(s) applied.
func NewEventManager(client *Client, opts ...EventManagerConfigOpt) EventManager {
	config := DefaultEventManagerConfig()
	config.Apply(opts)

	return &eventManagerImpl{
		client: client,
		config: *config,
	}
}

// EventManager lets you listen for specific events triggered by raw gateway events
type EventManager interface {
	// AddEventListeners adds one or more EventListener(s) to the EventManager
	AddEventListeners(eventListeners ...EventListener)

	// RemoveEventListeners removes one or more EventListener(s) from the EventManager
	RemoveEventListeners(eventListeners ...EventListener)

	// HandleEvent calls the correct EventListener(s) for the given gateway.Event
	HandleEvent(event gateway.Event)
}

// EventListener is used to create new EventListener to listen to events
type EventListener interface {
	OnEvent(c *Client, e gateway.Event)
}

// EventListenerFunc is a helper type to create a EventListener from a func(c *Client, e gateway.Event)
type EventListenerFunc func(c *Client, e gateway.Event)

func (f EventListenerFunc) OnEvent(c *Client, e gateway.Event) {
	f(c, e)
}

// NewListenerFunc returns a new EventListener for the given func(c *Client, e E)
func NewListenerFunc[E gateway.Event](f func(c *Client, e E)) EventListener {
	return &listenerFunc[E]{f: f}
}

type listenerFunc[E gateway.Event] struct {
	f func(c *Client, e E)
}

func (l *listenerFunc[E]) OnEvent(c *Client, e gateway.Event) {
	if event, ok := e.(E); ok {
		l.f(c, event)
	}
}

// NewListenerChan returns a new EventListener for the given chan<- gateway.Event
func NewListenerChan[E gateway.Event](c chan<- E) EventListener {
	return &listenerChan[E]{c: c}
}

type listenerChan[E gateway.Event] struct {
	c chan<- E
}

func (l *listenerChan[E]) OnEvent(_ *Client, e gateway.Event) {
	if event, ok := e.(E); ok {
		l.c <- event
	}
}

type eventManagerImpl struct {
	client          *Client
	eventListenerMu sync.RWMutex
	config          EventManagerConfig
}

func (m *eventManagerImpl) HandleEvent(event gateway.Event) {
	// set respond function if not set to handle http & gateway interactions the same way
	if e, ok := event.(gateway.EventInteractionCreate); ok && e.Respond == nil {
		e.Respond = func(response discord.InteractionResponse) error {
			return m.client.Rest.CreateInteractionResponse(e.Interaction.ID(), e.Interaction.Token(), response)
		}
		event = e
	}

	defer func() {
		if r := recover(); r != nil {
			m.config.Logger.Errorf("recovered from panic in event listener: %+v\nstack: %s", r, string(debug.Stack()))
			return
		}
	}()
	m.eventListenerMu.RLock()
	defer m.eventListenerMu.RUnlock()
	for i := range m.config.EventListeners {
		m.config.EventListeners[i].OnEvent(m.client, event)
	}
}

func (m *eventManagerImpl) AddEventListeners(listeners ...EventListener) {
	m.eventListenerMu.Lock()
	defer m.eventListenerMu.Unlock()
	m.config.EventListeners = append(m.config.EventListeners, listeners...)
}

func (m *eventManagerImpl) RemoveEventListeners(listeners ...EventListener) {
	m.eventListenerMu.Lock()
	defer m.eventListenerMu.Unlock()
	for _, listener := range listeners {
		for i, l := range m.config.EventListeners {
			if l == listener {
				m.config.EventListeners = append(m.config.EventListeners[:i], m.config.EventListeners[i+1:]...)
				break
			}
		}
	}
}
