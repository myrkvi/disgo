package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"

	"github.com/disgoorg/log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
)

var token = os.Getenv("disgo_token")

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Infof("disgo version: %s", disgo.Version)

	client, err := bot.New(token,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildMessages, gateway.IntentDirectMessages, gateway.IntentMessageContent)),
		bot.WithEventListenerFunc(onMessageCreate),
	)
	if err != nil {
		log.Fatal("error while building bot: ", err)
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func onMessageCreate(c *bot.Client, e gateway.EventMessageCreate) {
	if e.Message.Author.Bot || e.Message.Author.System {
		return
	}
	if e.Message.Content == "start" {
		go func() {
			ch, cls := bot.NewEventCollector(c, func(e2 gateway.EventMessageCreate) bool {
				return e.ChannelID == e2.ChannelID && e.Message.Author.ID == e2.Message.Author.ID && e2.Message.Content != ""
			})
			defer cls()
			i := 1
			str := ">>> "
			ctx, clsCtx := context.WithTimeout(context.Background(), 20*time.Second)
			defer clsCtx()
			for {
				select {
				case <-ctx.Done():
					_, _ = c.Rest.CreateMessage(e.ChannelID, discord.NewMessageCreateBuilder().SetContent("cancelled").Build())
					return

				case messageEvent := <-ch:
					str += strconv.Itoa(i) + ". " + messageEvent.Message.Content + "\n\n"

					if i == 3 {
						_, _ = c.Rest.CreateMessage(messageEvent.ChannelID, discord.NewMessageCreateBuilder().SetContent(str).Build())
						return
					}
					i++
				}
			}
		}()
	}
}
