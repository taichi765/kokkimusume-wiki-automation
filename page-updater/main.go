package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const API_ENDPOINT_BASE = "https://api.wikiwiki.jp/kokkimusume"

func main() {
	passwd, err := loadPassword()
	if err != nil {
		log.Fatalf("failed to load password: %v", err)
	}

	charas, err := loadCharaData()
	if err != nil {
		log.Fatalf("failed to load character data: %v", err)
	}

	tok := getAuthToken(passwd)
	fmt.Println("successfully got token")

	c := &http.Client{}

	charaListSrc := getPageContent(c, "キャラ一覧", tok)
	fmt.Println("successfully fetched content for character list page")

	menubarSrc := getPageContent(c, "MenuBar", tok)
	fmt.Println("successfully fetched content for menu bar")

	newCharaListSrc, err := editCharaListPage(charaListSrc, charas)
	if err != nil {
		log.Fatal(err)
	}
	newMenubarSrc, err := editMenuBar(menubarSrc, charas)
	if err != nil {
		log.Fatal(err)
	}

	if err := putPageContent(c, "キャラ一覧", newCharaListSrc, tok); err != nil {
		log.Fatal(err)
	}
	fmt.Println("successfully updated character list page")

	if err := putPageContent(c, "MenuBar", newMenubarSrc, tok); err != nil {
		log.Fatal(err)
	}
	fmt.Println("successfully updated menu bar")
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
func getAuthToken(passwd string) string {
	body, err := json.Marshal(AuthRequest{
		Password: passwd,
	})
	if err != nil {
		log.Fatal("failed to marshal auth request")
	}

	buf := bytes.NewBuffer(body)

	res, err := http.Post(API_ENDPOINT_BASE+"/auth", "application/json", buf)
	if err != nil {
		log.Fatal("failed to send request")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatal("something went wrong")
	}

	var resJson AuthResponse
	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		log.Fatal("failed to decode response")
	}

	return resJson.Token
}

// Gets page content from WikiWiki's REST API.
func getPageContent(c *http.Client, page string, tok string) string {
	endpoint := API_ENDPOINT_BASE + "/page/" + page
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		panic("request must be valid")
	}
	req.Header.Add("Authorization", "Bearer "+tok)

	res, err := c.Do(req)
	if err != nil {
		log.Fatal("failed to get page content")
	}
	defer res.Body.Close()

	var resJson GetPageResponse
	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		log.Fatal("failed to decode response")
	}

	return resJson.Source
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
