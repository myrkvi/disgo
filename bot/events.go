package bot

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/rest"
)

// gateway.EventType(s) for all the bot events
const (
	EventTypeApplicationCommandInteractionCreate gateway.EventType = "APPLICATION_COMMAND_INTERACTION_CREATE"
	EventTypeAutocompleteInteractionCreate       gateway.EventType = "AUTOCOMPLETE_INTERACTION_CREATE"
	EventTypeComponentInteractionCreate          gateway.EventType = "COMPONENT_INTERACTION_CREATE"
	EventTypeModalInteractionCreate              gateway.EventType = "MODAL_INTERACTION_CREATE"
)

// EventApplicationCommandInteractionCreate is called when a user uses an application command
type EventApplicationCommandInteractionCreate struct {
	discord.ApplicationCommandInteraction
	Respond gateway.RespondFunc `json:"-"`
}

// EventType returns EventTypeApplicationCommandInteractionCreate
func (EventApplicationCommandInteractionCreate) EventType() gateway.EventType {
	return EventTypeApplicationCommandInteractionCreate
}

// CreateMessage responds to the interaction with a new message.
func (e EventApplicationCommandInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e EventApplicationCommandInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// CreateModal responds to the interaction with a new modal.
func (e EventApplicationCommandInteractionCreate) CreateModal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeModal, modalCreate, opts...)
}

type EventAutocompleteInteractionCreate struct {
	discord.AutocompleteInteraction
	Respond gateway.RespondFunc `json:"-"`
}

// EventType returns EventTypeAutocompleteInteractionCreate
func (EventAutocompleteInteractionCreate) EventType() gateway.EventType {
	return EventTypeAutocompleteInteractionCreate
}

// Result responds to the interaction with a slice of choices.
func (e EventAutocompleteInteractionCreate) Result(choices []discord.AutocompleteChoice, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeApplicationCommandAutocompleteResult, discord.AutocompleteResult{Choices: choices}, opts...)
}

type EventComponentInteractionCreate struct {
	discord.ComponentInteraction
	Respond gateway.RespondFunc `json:"-"`
}

// EventType returns EventTypeComponentInteractionCreate
func (EventComponentInteractionCreate) EventType() gateway.EventType {
	return EventTypeComponentInteractionCreate
}

// CreateMessage responds to the interaction with a new message.
func (e EventComponentInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e EventComponentInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// UpdateMessage responds to the interaction with updating the message the component is from.
func (e EventComponentInteractionCreate) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeUpdateMessage, messageUpdate, opts...)
}

// DeferUpdateMessage responds to the interaction with nothing.
func (e EventComponentInteractionCreate) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeDeferredUpdateMessage, nil, opts...)
}

// CreateModal responds to the interaction with a new modal.
func (e EventComponentInteractionCreate) CreateModal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeModal, modalCreate, opts...)
}

// EventModalInteractionCreate is called when a user submits a modal
type EventModalInteractionCreate struct {
	discord.ModalInteraction
	Respond gateway.RespondFunc `json:"-"`
}

// EventType returns EventTypeModalInteractionCreate
func (EventModalInteractionCreate) EventType() gateway.EventType {
	return EventTypeModalInteractionCreate
}

// CreateMessage responds to the interaction with a new message.
func (e EventModalInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e EventModalInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// UpdateMessage responds to the interaction with updating the message the component is from.
func (e EventModalInteractionCreate) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeUpdateMessage, messageUpdate, opts...)
}

// DeferUpdateMessage responds to the interaction with nothing.
func (e EventModalInteractionCreate) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeDeferredUpdateMessage, nil, opts...)
}
