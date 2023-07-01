package oauth2

import (
	"github.com/disgoorg/disgo/internal/ttlmap"
	"github.com/disgoorg/log"
)

var _ StateController = (*defaultStateController)(nil)

// StateController is responsible for generating, storing and validating states.
type StateController interface {
	// NewState generates a new random state to be used as a state.
	NewState(redirectURI string) string

	// UseState validates a state and returns the redirect url or nil if it is invalid.
	UseState(state string) string
}

// NewStateController returns a new empty StateController.
func NewStateController(opts ...StateControllerConfigOpt) StateController {
	config := DefaultStateControllerConfig()
	config.Apply(opts)

	states := ttlmap.New(config.MaxTTL)
	for state, url := range config.States {
		states.Put(state, url)
	}

	return &defaultStateController{
		states:       states,
		newStateFunc: config.NewStateFunc,
		logger:       config.Logger,
	}
}

type defaultStateController struct {
	logger       log.Logger
	states       *ttlmap.Map
	newStateFunc func() string
}

func (c *defaultStateController) NewState(redirectURI string) string {
	state := c.newStateFunc()
	c.logger.Debugf("new state: %s for redirect uri: %s", state, redirectURI)
	c.states.Put(state, redirectURI)
	return state
}

func (c *defaultStateController) UseState(state string) string {
	uri := c.states.Get(state)
	if uri == "" {
		return ""
	}
	c.logger.Debugf("using state: %s for redirect uri: %s", state, uri)
	c.states.Delete(state)
	return uri
}
