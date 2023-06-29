package bot

import (
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/log"
)

// DefaultEventManagerConfig returns a new EventManagerConfig with all default values.
func DefaultEventManagerConfig() *EventManagerConfig {
	return &EventManagerConfig{
		Logger: log.Default(),
	}
}

// EventManagerConfig can be used to configure the EventManager.
type EventManagerConfig struct {
	Logger         log.Logger
	EventListeners []EventListener
}

// EventManagerConfigOpt is a functional option for configuring an EventManager.
type EventManagerConfigOpt func(config *EventManagerConfig)

// Apply applies the given EventManagerConfigOpt(s) to the EventManagerConfig.
func (c *EventManagerConfig) Apply(opts []EventManagerConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithEventManagerLogger overrides the default logger in the EventManagerConfig.
func WithEventManagerLogger(logger log.Logger) EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.Logger = logger
	}
}

// WithListeners adds the given EventListener(s) to the EventManagerConfig.
func WithListeners(listeners ...EventListener) EventManagerConfigOpt {
	return func(config *EventManagerConfig) {
		config.EventListeners = append(config.EventListeners, listeners...)
	}
}

// WithListenerFunc adds the given func(c *Client, e E) to the EventManagerConfig.
func WithListenerFunc[E gateway.Event](f func(c *Client, e E)) EventManagerConfigOpt {
	return WithListeners(NewListenerFunc(f))
}

// WithListenerChan adds the given chan<- E to the EventManagerConfig.
func WithListenerChan[E gateway.Event](c chan<- E) EventManagerConfigOpt {
	return WithListeners(NewListenerChan(c))
}
