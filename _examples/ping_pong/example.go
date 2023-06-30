package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/log"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, err := bot.New(os.Getenv("disgo_token"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
		bot.WithEventListenerFunc(onMessageCreate),
	)
	if err != nil {
		log.Fatal("error while building disgo: ", err)
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("errors while connecting to gateway: ", err)
	}

	log.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func onMessageCreate(c *bot.Client, e gateway.EventMessageCreate) {
	if e.Message.Author.Bot {
		return
	}
	var message string
	if e.Message.Content == "ping" {
		message = "pong"
	} else if e.Message.Content == "pong" {
		message = "ping"
	}
	if message != "" {
		_, _ = c.Rest.CreateMessage(e.ChannelID, discord.NewMessageCreateBuilder().SetContent(message).Build())
	}
}
