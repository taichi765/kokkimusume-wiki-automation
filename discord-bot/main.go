package discordbot

import (
	"fmt"
	"log"
	"os"

	"github.com/disgoorg/disgo"
	"github.com/joho/godotenv"
)

func main() {
	tok, err := loadDiscordToken()
	if err != nil {
		log.Fatalf("failed to load discord's token: %v", err)
	}
	client, err := disgo.New(tok)
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
