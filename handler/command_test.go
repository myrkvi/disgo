package handler_test

import (
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type SupportedFieldTypes struct {
	Member1    discord.Member
	Member2    discord.ResolvedMember
	Member3    discord.User
	Channel    discord.ResolvedChannel
	Role       discord.Role
	Attachment discord.Attachment
	Snowflake  snowflake.ID
	String     string
	Integer    int     // int32, int64, uint, uint32, uint64
	Float      float32 // float64
	Boolean    bool
}

var UserNoteCommand = discord.SlashCommandCreate{
	Name:        "note-add",
	Description: "Add a note about a user",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "target-user",
			Description: "The user that you want to make a note about.",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "note",
			Description: "The note that you want to add.",
			Required:    true,
		},
	},
}

type UserNoteArgs struct {
	TargetUser discord.Member `disgo:"target-user"` // Will try to get the parameter "target-user".
	Note       string         // Will by default try to get the parameter "note".
	Timestamp  uint           `disgo:"-"` // This field will be ignored
}

var e *handler.CommandEvent

func ExampleCommandEvent_ParseSlashCommandData() {
	var args UserNoteArgs
	e.ParseSlashCommandData(&args)
}
