package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const API_ENDPOINT_BASE = "https://api.wikiwiki.jp/kokkimusume"

func main() {
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

	tok, err := getAuthToken(passwd)
	if err != nil {
		slog.Error("failed to get token", slog.Any("err", err))
	}
	slog.Info("successfully got token")

	c := &http.Client{}

	charaListSrc, err := getPageContent(c, "キャラ一覧", tok)
	if err != nil {
		return 1
	}
	slog.Info("successfully fetched content for character list page")

	menubarSrc, err := getPageContent(c, "MenuBar", tok)
	if err != nil {
		return 1
	}
	slog.Info("successfully fetched content for menu bar")

	newCharaListSrc, err := editCharaListPage(charaListSrc, charas)
	if err != nil {
		slog.Error("failed to generate chara list page", slog.Any("err", err))
		return 1
	}
	newMenubarSrc, err := editMenuBar(menubarSrc, charas)
	if err != nil {
		slog.Error("failed to generate MenuBar", slog.Any("err", err))
		return 1
	}

	if err := putPageContent(c, "キャラ一覧", newCharaListSrc, tok); err != nil {
		slog.Error("failed to update chara list page", slog.Any("err", err))
		return 1
	}
	slog.Info("successfully updated character list page")

	if err := putPageContent(c, "MenuBar", newMenubarSrc, tok); err != nil {
		slog.Error("failed to update MenuBar", slog.Any("err", err))
		return 1
	}
	slog.Info("successfully updated menu bar")

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

// getAuthToken gets token from wikiwiki's REST API.
func getAuthToken(passwd string) (string, error) {
	body, err := json.Marshal(AuthRequest{
		Password: passwd,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal auth request: %w", err)
	}

	buf := bytes.NewBuffer(body)

	res, err := http.Post(API_ENDPOINT_BASE+"/auth", "application/json", buf)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("something went wrong with authentication")
	}

	var resJson AuthResponse
	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return resJson.Token, nil
}

// Gets page content from WikiWiki's REST API.
func getPageContent(c *http.Client, page string, tok string) (string, error) {
	endpoint := API_ENDPOINT_BASE + "/page/" + page
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		panic("request must be valid")
	}
	req.Header.Add("Authorization", "Bearer "+tok)

	res, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get page content: %w", err)
	}
	defer res.Body.Close()

	var resJson GetPageResponse
	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return resJson.Source, nil
}

// Updates page content using WikiWiki's REST API.
func putPageContent(c *http.Client, page string, content string, tok string) error {
	endpoint := API_ENDPOINT_BASE + "/page/" + page
	body, err := json.Marshal(PutPageRequest{
		Source: content,
	})
	if err != nil {
		panic("json must be valid")
	}

	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer(body))
	if err != nil {
		panic("request must be valid")
	}
	req.Header.Add("Authorization", "Bearer "+tok)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update content: %w", err)
	}
	defer res.Body.Close()

	var resJson PutPageResponse
	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if resJson.Status != "ok" {
		return fmt.Errorf("something went wrong while updating page content: status was %v", resJson.Status)
	}

	return nil
}
