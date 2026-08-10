package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/joho/godotenv"
)

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "hello",
		Description: "say hellp",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "name",
				Description: "your name",
				Required:    true,
			},
		},
	},
}

// 環境変数または.envからDicordのトークンを読む
func loadDiscordToken() (string, error) {
	tok, ok := os.LookupEnv("DISCORD_TOKEN")
	if ok {
		return tok, nil
	}

	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	tok, ok = os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		return "", fmt.Errorf("can't find DISCORD_TOKEN in both env vars and .env file")
	}
	return tok, nil
}

func main() {
	slog.Info("loading discord token")
	tok, err := loadDiscordToken()

	if err != nil {
		slog.Error("failed to load discord's token", slog.Any("err", err))
	}
	slog.Info("successfully loaded discord token")

	client, err := disgo.New(tok)
	if err != nil {
		log.Fatalf("failed to build client")
	}
	defer client.Close(context.TODO())

	_, err = client.Rest.SetGlobalCommands(client.ApplicationID, commands)
	if err != nil {
		log.Fatalf("failed to register commands: %v", err)
	}
	fmt.Println("successfully registered application commands")
}
