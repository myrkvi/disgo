package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/log"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/sharding"
)

var token = os.Getenv("disgo_token")

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Info("disgo version: ", disgo.Version)

	client, err := bot.New(token,
		bot.WithShardManagerConfigOpts(
			sharding.WithShardIDs(0, 1),
			sharding.WithShardCount(2),
			sharding.WithAutoScaling(true),
			sharding.WithGatewayConfigOpts(
				gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildMessages, gateway.IntentDirectMessages),
				gateway.WithCompress(true),
			),
		),
		bot.WithEventListeners(
			bot.NewListenerFunc(onMessageCreate),
			bot.NewListenerFunc(func(c *bot.Client, e gateway.EventReady) {
				log.Infof("shard [%d/%d] ready", e.Shard[0], e.Shard[1])
			}),
		),
	)
	if err != nil {
		log.Fatalf("error while building disgo: %s", err)
	}

	defer client.Close(context.TODO())

	if err = client.OpenShardManager(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func onMessageCreate(c *bot.Client, e gateway.EventMessageCreate) {
	if e.Message.Author.Bot {
		return
	}
	_, _ = c.Rest.CreateMessage(e.ChannelID, discord.NewMessageCreateBuilder().SetContent(e.Message.Content).Build())
}
