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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to read dotenv")
	}

	tok := getAuthToken()
	fmt.Println(tok)
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

	res, err := http.Post("https://api.wikiwiki.jp/kokkimusume/auth", "application/json", buf)
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
