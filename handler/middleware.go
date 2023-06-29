package handler

import (
	"github.com/disgoorg/disgo/gateway"
)

type (
	Handler func(e gateway.EventInteractionCreate) error

	Middleware func(next Handler) Handler

	Middlewares []Middleware
)
