package bot

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

const (
	EventTypeApplicationCommandInteractionCreate gateway.EventType = "APPLICATION_COMMAND_INTERACTION_CREATE"
	EventTypeAutocompleteInteractionCreate       gateway.EventType = "AUTOCOMPLETE_INTERACTION_CREATE"
	EventTypeComponentInteractionCreate          gateway.EventType = "COMPONENT_INTERACTION_CREATE"
	EventTypeModalInteractionCreate              gateway.EventType = "MODAL_INTERACTION_CREATE"
)

type EventApplicationCommandInteractionCreate struct {
	discord.ApplicationCommandInteraction
	Respond gateway.RespondFunc `json:"-"`
}

func (EventApplicationCommandInteractionCreate) EventType() gateway.EventType {
	return EventTypeApplicationCommandInteractionCreate
}

type EventAutocompleteInteractionCreate struct {
	discord.AutocompleteInteraction
	Respond gateway.RespondFunc `json:"-"`
}

func (EventAutocompleteInteractionCreate) EventType() gateway.EventType {
	return EventTypeAutocompleteInteractionCreate
}

type EventComponentInteractionCreate struct {
	discord.ComponentInteraction
	Respond gateway.RespondFunc `json:"-"`
}

func (EventComponentInteractionCreate) EventType() gateway.EventType {
	return EventTypeComponentInteractionCreate
}

type EventModalInteractionCreate struct {
	discord.ModalInteraction
	Respond gateway.RespondFunc `json:"-"`
}

func (EventModalInteractionCreate) EventType() gateway.EventType {
	return EventTypeModalInteractionCreate
}
