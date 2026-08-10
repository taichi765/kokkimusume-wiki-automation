package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v88/github"
	"github.com/taichi765/kokkimusume-wiki-automation/common"
)

const REPO_OWNER = "taichi765"
const REPO_NAME = "kokkimusume-wiki-automation"
const CHARA_CSV_PATH = "data/chara.csv"

func updateCsv(chara common.CharacterData, appId, installationId int64) error {
	client, err := newClient(appId, installationId)
	if err != nil {
		return err
	}

	content, _, _, err := client.Repositories.GetContents(context.TODO(), REPO_OWNER, REPO_NAME, CHARA_CSV_PATH, &github.RepositoryContentGetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get file content: %w", err)
	}

	sha := content.GetSHA()
	old, err := content.GetContent()
	if err != nil {
		return fmt.Errorf("failed to read content from reponse: %w", err)
	}

	b := bytes.NewBuffer([]byte(old))
	wr := csv.NewWriter(b)

	if err := wr.Write([]string{chara.Name, chara.Area, chara.FirstAppearenceDate}); err != nil {
		return fmt.Errorf("failed to append line to CSV file: %w", err)
	}
	wr.Flush()
	if wr.Error() != nil {
		return fmt.Errorf("failed to flush csv content")
	}

	msg := "Update chara.csv"
	_, _, err = client.Repositories.UpdateFile(context.TODO(), REPO_OWNER, REPO_NAME, CHARA_CSV_PATH, &github.RepositoryContentFileOptions{
		Message: &msg,
		Content: b.Bytes(),
		SHA:     &sha,
	})
	if err != nil {
		return fmt.Errorf("failed to update csv file: %w", err)
	}

	return nil
}

func newClient(appId, installationId int64) (*github.Client, error) {
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appId, installationId, "./kokkimusume-wiki-automation.2026-08-10.private-key.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticated transport: %w", err)
	}

	client, err := github.NewClient(github.WithTransport(itr))
	if err != nil {
		return nil, fmt.Errorf("failed to build github client: %w", err)
	}

	return client, nil
}
