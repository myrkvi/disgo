package bot

import (
	"runtime/debug"
	"sync"

	"github.com/disgoorg/disgo/gateway"
)

var _ EventManager = (*eventManagerImpl)(nil)

// NewEventManager returns a new EventManager with the EventManagerConfigOpt(s) applied.
func NewEventManager(client Client, opts ...EventManagerConfigOpt) EventManager {
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
	OnEvent(event gateway.Event)
}

// NewListenerFunc returns a new EventListener for the given func(e E)
func NewListenerFunc[E gateway.Event](f func(e E)) EventListener {
	return &listenerFunc[E]{f: f}
}

type listenerFunc[E gateway.Event] struct {
	f func(e E)
}

func (l *listenerFunc[E]) OnEvent(e gateway.Event) {
	if event, ok := e.(E); ok {
		l.f(event)
	}
}

// NewListenerChan returns a new EventListener for the given chan<- gateway.Event
func NewListenerChan[E gateway.Event](c chan<- E) EventListener {
	return &listenerChan[E]{c: c}
}

type listenerChan[E gateway.Event] struct {
	c chan<- E
}

func (l *listenerChan[E]) OnEvent(e gateway.Event) {
	if event, ok := e.(E); ok {
		l.c <- event
	}
}

type eventManagerImpl struct {
	client          Client
	eventListenerMu sync.Mutex
	config          EventManagerConfig
}

func (e *eventManagerImpl) HandleEvent(event gateway.Event) {
	defer func() {
		if r := recover(); r != nil {
			e.config.Logger.Errorf("recovered from panic in event listener: %+v\nstack: %s", r, string(debug.Stack()))
			return
		}
	}()
	e.eventListenerMu.Lock()
	defer e.eventListenerMu.Unlock()
	for i := range e.config.EventListeners {
		e.config.EventListeners[i].OnEvent(event)
	}
}

func (e *eventManagerImpl) AddEventListeners(listeners ...EventListener) {
	e.eventListenerMu.Lock()
	defer e.eventListenerMu.Unlock()
	e.config.EventListeners = append(e.config.EventListeners, listeners...)
}

func (e *eventManagerImpl) RemoveEventListeners(listeners ...EventListener) {
	e.eventListenerMu.Lock()
	defer e.eventListenerMu.Unlock()
	for _, listener := range listeners {
		for i, l := range e.config.EventListeners {
			if l == listener {
				e.config.EventListeners = append(e.config.EventListeners[:i], e.config.EventListeners[i+1:]...)
				break
			}
		}
	}
}
