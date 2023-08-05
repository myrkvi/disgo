package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = snowflake.GetEnv("disgo_guild_id")

	//go:embed gopher.png
	gopher []byte
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Info("bot version: ", disgo.Version)

	client, err := bot.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentsNonPrivileged, gateway.IntentMessageContent),
			gateway.WithPresenceOpts(gateway.WithListeningActivity("your bullshit"), gateway.WithOnlineStatus(discord.OnlineStatusDND)),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagsAll),
		),
		bot.WithMemberChunkingFilter(bot.MemberChunkingFilterNone),
		bot.WithEventListeners(
			bot.NewListenerFunc(onMessageUpdate),
			bot.NewListenerFunc(onInteraction),
		),
	)
	if err != nil {
		log.Fatal("error while building bot instance: ", err)
		return
	}

	registerCommands(client)

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to discord: ", err)
	}

	defer client.Close(context.TODO())

	log.Info("ExampleBot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func onMessageUpdate(c *bot.Client, e gateway.EventMessageUpdate) {
	fmt.Printf("old message: %+v\n", e.OldMessage)
	fmt.Printf("new message: %+v\n", e.Message)
}

func onInteraction(c *bot.Client, e gateway.EventInteractionCreate) {
	switch e.Interaction.(type) {
	case discord.ApplicationCommandInteraction:
		_ = e.Respond(discord.InteractionResponseTypeCreateMessage, discord.MessageCreate{
			Content: "Hello World!",
		})
	case discord.ComponentInteraction:

	case discord.ModalInteraction:

	}
}
