package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/taichi765/kokkimusume-wiki-automation/wikiwiki"
)

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	slog.SetDefault(l)
	os.Exit(runMain())
}

func runMain() int {
	passwd, err := loadPassword()
	if err != nil {
		slog.Error("failed to load password", slog.Any("err", err))
		return 1
	}

	charas, err := loadCharaData()
	if err != nil {
		slog.Error("failed to load character data", slog.Any("err", err))
		return 1
	}

	tok, err := wikiwiki.GetAuthToken(passwd)
	if err != nil {
		slog.Error("failed to get token", slog.Any("err", err))
		return 1
	}
	slog.Info("successfully got token")

	c := &http.Client{}

	charaListSrc, err := wikiwiki.FetchPageContent(c, "キャラ一覧", tok)
	if err != nil {
		slog.Error("failed to fetch character list page content", slog.Any("err", err))
		return 1
	}

	menubarSrc, err := wikiwiki.FetchPageContent(c, "MenuBar", tok)
	if err != nil {
		slog.Error("failed to fetch MenuBar content", slog.Any("err", err))
		return 1
	}

	newCharaListSrc, err := generateCharaListPage(charaListSrc, charas)
	if err != nil {
		slog.Error("failed to generate chara list page", slog.Any("err", err))
		return 1
	}
	newMenubarSrc, err := generateMenuBar(menubarSrc, charas)
	if err != nil {
		slog.Error("failed to generate MenuBar", slog.Any("err", err))
		return 1
	}

	if err := wikiwiki.UpdatePageContent(c, "キャラ一覧", newCharaListSrc, tok); err != nil {
		slog.Error("failed to update chara list page", slog.Any("err", err))
		return 1
	}

	if err := wikiwiki.UpdatePageContent(c, "MenuBar", newMenubarSrc, tok); err != nil {
		slog.Error("failed to update MenuBar", slog.Any("err", err))
		return 1
	}

	return 0
}

// 環境変数または.envからパスワードを取得する
func loadPassword() (string, error) {
	passwd, ok := os.LookupEnv("WIKIWIKI_PASSWORD")
	if ok {
		return passwd, nil
	}

	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("failed to read dotenv")
	}

	passwd, ok = os.LookupEnv("WIKIWIKI_PASSWORD")
	if !ok {
		return "", fmt.Errorf("can't find WIKIWIKI_PASSWORD in both env vars and .env file")
	}
	return passwd, nil
}
