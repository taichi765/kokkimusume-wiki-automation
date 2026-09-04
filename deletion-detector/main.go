package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
	"github.com/taichi765/kokkimusume-wiki-automation/wikiwiki"
)

const azureBlobStorageUrl = "https://kokkimusumedetector.blob.core.windows.net"
const azureBlobContainerName = "page-list"

const deletedThreshould = 5

const discordNotifyChanel = "1537321425944191006"

type EnvVars struct {
	discordToken     string
	wikiwikiPassword string
}

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	slog.SetDefault(l)
	os.Exit(runMain())
}

func runMain() int {
	envVars, err := loadEnvVars()
	if err != nil {
		slog.Error("failed to load env vars", slog.Any("err", err))
		return 1
	}

	// WikiWiki
	wc, err := wikiwiki.NewClient(envVars.wikiwikiPassword)
	if err != nil {
		slog.Error("failed to create wikiwiki client", slog.Any("err", err))
		return 1
	}
	slog.Debug("succeed to create wikiwiki client")

	curr, err := wc.GetPageList()
	if err != nil {
		slog.Error("failed to get the list of pages", slog.Any("err", err))
		return 1
	}
	slog.Info("successfully got the list of pages")

	// Azure Blob Storage
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		slog.Error("failed to authenticate on Azure", slog.Any("err", err))
		return 1
	}

	client, err := azblob.NewClient(azureBlobStorageUrl, cred, nil)
	if err != nil {
		slog.Error("failed to create azure blob storage client", slog.Any("err", err))
		return 1
	}

	prev, err := getPreviousPageList(client)
	if err != nil {
		slog.Error("failed to get previous page list", slog.Any("err", err))
		return 1
	}

	// TODO: createdでも荒らし判定するか？
	_, deleted := countPageChanges(prev.Pages, curr.Pages)
	if deleted >= deletedThreshould {
		err := notifyViaDiscord(envVars.discordToken, deleted)
		if err != nil {
			slog.Error("failed to notify via discord", slog.Any("err", err))
			return 1
		}
	}

	if err := uploadPageList(client, curr); err != nil {
		slog.Error("failed to upload latest page list", slog.Any("err", err))
		return 1
	}

	return 0
}

// loadEnvVars reads secrets and variables from environment variable or .env file.
func loadEnvVars() (*EnvVars, error) {
	dotEnvLoaded := false
	load := func(name string) (string, error) {
		v, ok := os.LookupEnv(name)
		if !ok && dotEnvLoaded {
			return "", fmt.Errorf("can't find %s in both env vars and .env file", name)
		}
		if !ok && !dotEnvLoaded {
			err := godotenv.Load()
			if err != nil {
				return "", fmt.Errorf("failed to load dotenv: %w", err)
			}
			dotEnvLoaded = true

			v, ok = os.LookupEnv(name)
			if !ok {
				return "", fmt.Errorf("can't find %s in both env vars and .env file", name)
			}
		}
		return v, nil
	}

	passwd, err := load("WIKIWIKI_PASSWORD")
	if err != nil {
		return nil, err
	}

	tok, err := load("DISCORD_TOKEN")
	if err != nil {
		return nil, err
	}

	return &EnvVars{
		wikiwikiPassword: passwd,
		discordToken:     tok,
	}, nil
}

func getPreviousPageList(c *azblob.Client) (wikiwiki.GetPageListResponse, error) {
	pager := c.NewListBlobsFlatPager(azureBlobContainerName, nil)
	var times []time.Time

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			return wikiwiki.GetPageListResponse{}, fmt.Errorf("an error occureed while finding previous page list: %w", err)
		}

		for _, blob := range resp.Segment.BlobItems {
			time, err := time.Parse(time.DateTime, *blob.Name)
			if err != nil {
				return wikiwiki.GetPageListResponse{}, fmt.Errorf("blob name was invalid: %w", err)
			}
			times = append(times, time)
		}
	}

	slices.SortFunc(times, func(a, b time.Time) int {
		return a.Compare(b)
	})

	blobName := times[len(times)-1].Format(time.DateTime)

	resp, err := c.DownloadStream(context.TODO(), azureBlobContainerName, blobName, nil)
	if err != nil {
		return wikiwiki.GetPageListResponse{}, fmt.Errorf("failed to download previous page list: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}()

	var data wikiwiki.GetPageListResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return wikiwiki.GetPageListResponse{}, fmt.Errorf("failed to decode prevopus page list: %w", err)
	}

	return data, nil
}

func uploadPageList(c *azblob.Client, pages wikiwiki.GetPageListResponse) error {
	blobName := time.Now().Format(time.DateTime)
	data, err := json.Marshal(pages)
	if err != nil {
		return fmt.Errorf("failed to encode page list: %w", err)
	}

	_, err = c.UploadBuffer(context.TODO(), azureBlobContainerName, blobName, data, nil)
	if err != nil {
		return fmt.Errorf("failed to upload page list: %w", err)
	}

	return nil
}

// countPageChanges counts number of created pages and deleted pages.
//
// First return value is number of created pages, and secound is of deleted pages.
func countPageChanges(prev, curr []wikiwiki.GeneralPageInfo) (int, int) {
	names := make(map[string]bool)

	for _, p := range prev {
		names[p.Name] = false
	}

	var created []string
	var deleted []string

	for _, p := range curr {
		// prevに存在しないがcurrに存在する
		if _, ok := names[p.Name]; !ok {
			created = append(created, p.Name)
		} else {
			names[p.Name] = true
		}
	}

	for name, exists := range names {
		// prevに存在するがcurrに存在しない
		if !exists {
			deleted = append(deleted, name)
		}
	}

	return len(created), len(deleted)
}

func notifyViaDiscord(tok string, deleted int) error {
	client, err := disgo.New(tok)
	if err != nil {
		return fmt.Errorf("error while building disgo: %w", err)
	}

	channel := snowflake.MustParse(discordNotifyChanel)

	_, err = client.Rest.CreateMessage(channel, discord.NewMessageCreate().WithContentf("%vページが削除されました", deleted))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
