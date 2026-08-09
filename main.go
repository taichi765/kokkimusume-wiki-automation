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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to read dotenv")
	}

	tok := getAuthToken()
	c := http.Client{}
	src := getPageContent(&c, "FrontPage", tok)
	fmt.Println(src)
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
