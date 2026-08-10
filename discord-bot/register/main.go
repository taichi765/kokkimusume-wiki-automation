package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/taichi765/kokkimusume-wiki-automation/discord-bot/utils"
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

func main() {
	slog.Info("loading discord token")
	tok, err := utils.LoadDiscordToken()

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
