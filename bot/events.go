package bot

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

type EventApplicationCommandInteractionCreate struct {
	discord.ApplicationCommandInteraction
	Respond gateway.RespondFunc `json:"-"`
}

type EventAutocompleteInteractionCreate struct {
	discord.AutocompleteInteraction
	Respond gateway.RespondFunc `json:"-"`
}

type EventComponentInteractionCreate struct {
	discord.ComponentInteraction
	Respond gateway.RespondFunc `json:"-"`
}

type EventModalInteractionCreate struct {
	discord.ModalInteraction
	Respond gateway.RespondFunc `json:"-"`
}
