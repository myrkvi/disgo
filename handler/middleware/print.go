package middleware

import (
	"github.com/disgoorg/disgo/handler"
)

func Print(content string) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			println(content)
			return next(e)
		}
	}
}
