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
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
)

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "new",
		Description: "append new character to CSV file on github",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleJapanese: "新しい国旗娘をGitHub上のCSVファイルに追加する",
		},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall,
		},
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
		},
	},
	discord.SlashCommandCreate{
		Name:        "dummy",
		Description: "dummy slash command",
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall,
		},
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeGuild,
		},
	},
}

type App struct {
	envVars *EnvVars
	client  *bot.Client
}

type EnvVars struct {
	githubAppId          int64
	githubInstallationId int64
	githubPrivateKey     []byte
	discordAppId         string
	discordToken         string
	discordPublicKey     string
}

// 環境変数または.envから[EnvVars]を読み込む
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
	githubPrivateKey, err := load("GITHUB_PRIVATE_KEY")
	if err != nil {
		return nil, err
	}

	return &EnvVars{
		discordToken:         discordToken,
		discordAppId:         discordAppId,
		discordPublicKey:     discordPublicKey,
		githubAppId:          githubAppId,
		githubInstallationId: githubInstallationId,
		githubPrivateKey:     []byte(githubPrivateKey),
	}, nil
}

func main() {
	os.Exit(runMain())
}

func runMain() int {
	slog.Info("loading env vars")
	envVars, err := loadEnvVars()
	if err != nil {
		slog.Error("failed to load env vars", slog.Any("err", err))
		return 1
	}
	slog.Info("successfully loaded env vars")

	app := App{
		envVars: envVars,
		client:  nil,
	}

	h := handler.New()
	h.SlashCommand("/new", newCharaModalSlashCommand)
	h.Modal("/modals/new", app.onNewCharaModalSubmitted)

	tok := envVars.discordToken
	client, err := disgo.New(tok,
		bot.WithHTTPServerConfigOpts(app.envVars.discordPublicKey,
			httpserver.WithURL("/interactions/"),
			httpserver.WithAddress(":8080"),
		),
		bot.WithEventListeners(h),
	)
	if err != nil {
		slog.Error("error while building disgo", slog.Any("err", err))
		return 1
	}
	defer client.Close(context.TODO())
	app.client = client

	devGuildId, ok := os.LookupEnv("DEV_DISCORD_GUILD_ID")
	guildIds := []snowflake.ID{}
	if ok {
		id, err := snowflake.Parse(devGuildId)
		if err != nil {
			slog.Error("can't parse DEV_DISCORD_GUILD_ID", slog.Any("err", err))
			return 1
		}
		guildIds = append(guildIds, id)
	}
	if err := handler.SyncCommands(client, commands, guildIds); err != nil {
		slog.Error("failed to register commands", slog.Any("err", err))
	}

	slog.Info("opening HTTP server")
	if err = client.OpenHTTPServer(); err != nil {
		slog.Error("failed to open HTTP server", slog.Any("err", err))
		return 1
	}

	slog.Info("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	return 0
}
