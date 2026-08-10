package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const API_ENDPOINT_BASE = "https://api.wikiwiki.jp/kokkimusume"

type CharacterData struct {
	Name                string
	Area                string
	FirstAppearenceDate string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to read dotenv")
	}

	charas, err := loadCharaData()
	if err != nil {
		log.Fatalf("failed to load character data: %v", err)
	}

	tok := getAuthToken()
	c := &http.Client{}
	charaListSrc := getPageContent(c, "キャラ一覧", tok)
	menubarSrc := getPageContent(c, "MenuBar", tok)

	newCharaListSrc, err := editCharaListPage(charaListSrc, charas)
	if err != nil {
		log.Fatal(err)
	}
	newMenubarSrc, err := editMenuBar(menubarSrc)
	if err != nil {
		log.Fatal(err)
	}

	if err := putPageContent(c, "キャラ一覧", newCharaListSrc, tok); err != nil {
		log.Fatal(err)
	}
	if err := putPageContent(c, "MenuBar", newMenubarSrc, tok); err != nil {
		log.Fatal(err)
	}
}

// getAuthToken gets token from wikiwiki's REST API.
func getAuthToken() string {
	body, err := json.Marshal(AuthRequest{
		Password: os.Getenv("WIKIWIKI_PASSWORD"),
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
	panic("TODO")
}
