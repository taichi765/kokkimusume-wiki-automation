package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/joho/godotenv"
)

type EnvVars struct {
	githubAppId          int64
	githubInstallationId int64
	discordAppId         string
	discordToken         string
	discordPublicKey     string
}

func loadEnvVars() (*EnvVars, error) {
	dotEnvLoaded := false
	load := func(name string) (string, error) {
		v, ok := os.LookupEnv(name)
		if !ok && dotEnvLoaded {
			return "", fmt.Errorf("can't find %s in both env vars and .env file", name)
		}
		if !ok && !dotEnvLoaded {
			err := godotenv.Load()
			if err != nil {
				return "", fmt.Errorf("failed to load dotenv: %w", err)
			}
			dotEnvLoaded = true

			v, ok = os.LookupEnv(name)
			if !ok {
				return "", fmt.Errorf("can't find %s in both env vars and .env file", name)
			}
		}
		return v, nil
	}

	discordToken, err := load("DISCORD_TOKEN")
	if err != nil {
		return nil, err
	}
	discordAppId, err := load("DISCORD_APP_ID")
	if err != nil {
		return nil, err
	}
	discordPublicKey, err := load("DISCORD_PUBLIC_KEY")
	if err != nil {
		return nil, err
	}
	githubAppIdStr, err := load("GITHUB_APP_ID")
	if err != nil {
		return nil, err
	}
	githubAppId, err := strconv.ParseInt(githubAppIdStr, 10, 64)
	if err != nil {
		return nil, err
	}
	githubInstallationIdStr, err := load("GITHUB_INSTALLATION_ID")
	if err != nil {
		return nil, err
	}
	githubInstallationId, err := strconv.ParseInt(githubInstallationIdStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return &EnvVars{
		discordToken:         discordToken,
		discordAppId:         discordAppId,
		discordPublicKey:     discordPublicKey,
		githubAppId:          githubAppId,
		githubInstallationId: githubInstallationId,
	}, nil
}

func main() {
	slog.Info("loading env vars")
	envVars, err := loadEnvVars()
	if err != nil {
		slog.Error("failed to load env vars", slog.Any("err", err))
	}
	slog.Info("successfully loaded env vars")

	client, err := disgo.New(envVars.discordToken,
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
