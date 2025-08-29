package handler

import (
	"context"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/internal/optparse"
	"github.com/disgoorg/disgo/rest"
)

// CommandEvent allows to handle all types of application command interactions.
type CommandEvent struct {
	*events.ApplicationCommandInteractionCreate
	Vars map[string]string
	Ctx  context.Context
}

func (e *CommandEvent) GetInteractionResponse(opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.GetInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *CommandEvent) UpdateInteractionResponse(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), messageUpdate, opts...)
}

func (e *CommandEvent) DeleteInteractionResponse(opts ...rest.RequestOpt) error {
	return e.Client().Rest.DeleteInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *CommandEvent) GetFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.GetFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}

func (e *CommandEvent) CreateFollowupMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), messageCreate, opts...)
}

func (e *CommandEvent) UpdateFollowupMessage(messageID snowflake.ID, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.UpdateFollowupMessage(e.ApplicationID(), e.Token(), messageID, messageUpdate, opts...)
}

func (e *CommandEvent) DeleteFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) error {
	return e.Client().Rest.DeleteFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}

// ParseSlashCommandData takes a reference to any struct, and tries to parse
// slash command interaction data into it using reflection. By default, it will
// try to match the struct's name in lowercase, but you can override it with
// the following struct tag: `disgo:"my-argument"`. You can prevent fields from
// being matched using `disgo:"-"`.
//
// Supports discord.ResolvedMember, discord.Member, and discord.User for
// user arguments. Supports discord.ResolvedChannel for channels.
// Supports other types for other argument types, including primitives.
//
// See the example for more details.
func (e *CommandEvent) ParseSlashCommandData(structRef any) {
	optparse.ParseCommandArguments(e.SlashCommandInteractionData(), structRef)
}
