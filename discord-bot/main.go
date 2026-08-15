package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
)

// set in ldFlags
var version = "unknown"
var commitHash = "unknown"

type App struct {
	envVars   *EnvVars
	isStarted *atomic.Bool
	isReady   *atomic.Bool
	isLiving  *atomic.Bool
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
	app := App{
		envVars:   nil,
		isStarted: &atomic.Bool{},
		isReady:   &atomic.Bool{},
		isLiving:  &atomic.Bool{},
	}
	app.isLiving.Store(true)
	app.openACAProbeServer()

	slog.Info("loading env vars")
	envVars, err := loadEnvVars()
	if err != nil {
		slog.Error("failed to load env vars", slog.Any("err", err))
		return 1
	}
	app.envVars = envVars
	slog.Info("successfully loaded env vars")

	code := app.openDiscordServer()
	if code != 0 {
		return code
	}
	app.isStarted.Store(true)
	app.isReady.Store(true)

	slog.Info("bot server is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	return 0
}

func (a *App) openDiscordServer() int {
	h := handler.New()
	h.SlashCommand("/new", newCharaModalSlashCommand)
	h.SlashCommand("/version", versionSlashCommand)
	h.SlashCommand("/help", helpSlashCommand)
	h.Modal("/modals/new", a.onNewCharaModalSubmitted)

	tok := a.envVars.discordToken
	client, err := disgo.New(tok,
		bot.WithHTTPServerConfigOpts(a.envVars.discordPublicKey,
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

	return 0
}

// Provides endpoint for health probe in Azure Container Apps.
func (a *App) openACAProbeServer() {
	// TODO: test that this function doesn't use a.envVars,
	// since it's called before startup phase completes.

	addr := ":8081"

	slog.Info("starting ACA healthy probe server", slog.String("addr", addr))

	mux := a.newACAProbeServeMux()
	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	go func() {
		slog.Info("startup probe server is now listening...")
		err := s.ListenAndServe()
		if err != nil {
			slog.Error("something went wrong in startup server", slog.Any("err", err))
		}
	}()
}

// Creates [http.ServeMux] to respond to ACA's health probes.
func (a *App) newACAProbeServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/startup", a.startupHandler)
	mux.HandleFunc("/readiness", a.readinessHandler)
	mux.HandleFunc("/liveness", a.livenessHandler)
	return mux
}

// Used with [http.HandleFunc] in [newACAProbeServer]
func (a *App) startupHandler(w http.ResponseWriter, r *http.Request) {
	if a.isStarted.Load() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (a *App) readinessHandler(w http.ResponseWriter, r *http.Request) {
	if a.isReady.Load() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (a *App) livenessHandler(w http.ResponseWriter, r *http.Request) {
	if a.isLiving.Load() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
