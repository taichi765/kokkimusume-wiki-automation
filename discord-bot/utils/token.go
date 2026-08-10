package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// 環境変数または.envからDicordのトークンを読む
func LoadDiscordToken() (string, error) {
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
