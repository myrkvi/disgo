package handler

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
)

type (
	InteractionEvent struct {
		gateway.EventInteractionCreate
		Client *bot.Client
		Vars   map[string]string
	}

	Handler func(e *InteractionEvent) error

	Middleware func(next Handler) Handler

	Middlewares []Middleware
)
