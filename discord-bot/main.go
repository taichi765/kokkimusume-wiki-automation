package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/taichi765/kokkimusume-wiki-automation/discord-bot/utils"
)

func main() {
	slog.Info("loading discord token")
	tok, err := utils.LoadDiscordToken()
	if err != nil {
		slog.Error("failed to load discord's token", slog.Any("err", err))
	}
	slog.Info("successfully loaded discord token")

	client, err := disgo.New(tok,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMembers,
				gateway.IntentMessageContent,
			),
		),
		bot.WithEventListenerFunc(commandListener))
	if err != nil {
		slog.Error("error while building disgo", slog.Any("err", err))
		return
	}

	defer client.Close(context.TODO())

	slog.Info("opening gateway between discord")
	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("errors while connecting to gateway", slog.Any("err", err))
		return
	}

	slog.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func commandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	switch data.CommandName() {
	case "hello":
		err := event.CreateMessage(discord.NewMessageCreate().WithContentf("Hello, %v!", data.String("name")))
		if err != nil {
			slog.Error("failed to send message", slog.Any("err", err))
		}
	}
}
