package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = snowflake.GetEnv("disgo_guild_id")
)

func main() {
	log.SetLevel(log.LevelInfo)
	log.Info("starting example...")
	log.Info("disgo version: ", disgo.Version)

	client, err := bot.New(token,
		bot.WithDefaultGateway(),
		bot.WithEventListenerFunc(eventListenerFunc),
		bot.WithEventListeners(&eventListener{}),
	)
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
		return
	}
	client.AddEventListeners(bot.NewListenerChan(eventListenerChan(client)))

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

type eventListener struct{}

func (*eventListener) OnEvent(c *bot.Client, e gateway.Event) {
	switch event := e.(type) {
	case gateway.EventMessageCreate:
		if event.Message.Content == "ping" {
			_, _ = c.Rest.CreateMessage(event.ChannelID, discord.MessageCreate{
				Content: "pong",
			})
		}
	}
}

func eventListenerFunc(c *bot.Client, e gateway.EventMessageCreate) {
	_, _ = c.Rest.CreateMessage(e.ChannelID, discord.MessageCreate{
		Content: "pong",
	})
}

func eventListenerChan(c *bot.Client) chan<- gateway.EventMessageCreate {
	events := make(chan gateway.EventMessageCreate)
	go func() {
		defer close(events)
		for e := range events {
			if e.Message.Content == "ping" {
				_, _ = c.Rest.CreateMessage(e.ChannelID, discord.MessageCreate{
					Content: "pong",
				})
			}
		}
	}()
	return events
}
