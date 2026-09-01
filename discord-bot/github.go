package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-github/v88/github"
)

const REPO_OWNER = "taichi765"
const REPO_NAME = "kokkimusume-wiki-automation"
const CHARA_CSV_PATH = "data/charas.csv"

// Error returned from [updateCsv].
type CharaAlreadyExistsError struct {
	Name string
}

func (e *CharaAlreadyExistsError) Error() string {
	return fmt.Sprintf("charater %s already exists", e.Name)
}

func updateCsv(chara common.CharacterData, c *github.Client) error {
	content, err := getCsvContent(c)
	if err != nil {
		return err
	}

	sha := content.GetSHA()
	old, err := content.GetContent()
	if err != nil {
		return fmt.Errorf("failed to read content from reponse: %w", err)
	}

	rd := csv.NewReader(strings.NewReader(old))
	for {
		r, err := rd.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV file: %w", err)
		}

		if r[0] == chara.Name {
			return &CharaAlreadyExistsError{chara.Name}
		}
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
	_, _, err = c.Repositories.UpdateFile(context.TODO(), REPO_OWNER, REPO_NAME, CHARA_CSV_PATH, &github.RepositoryContentFileOptions{
		Message: &msg,
		Content: b.Bytes(),
		SHA:     &sha,
	})
	if err != nil {
		return fmt.Errorf("failed to update csv file: %w", err)
	}

	return nil
}

func getCsvContent(c *github.Client) (*github.RepositoryContent, error) {
	content, _, _, err := c.Repositories.GetContents(context.TODO(), REPO_OWNER, REPO_NAME, CHARA_CSV_PATH, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file content: %w", err)
	}
	return content, nil
}

func newClient(appId, installationId int64, privateKey []byte) (*github.Client, error) {
	itr, err := ghinstallation.New(http.DefaultTransport, appId, installationId, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticated transport: %w", err)
	}

	client, err := github.NewClient(github.WithTransport(itr))
	if err != nil {
		return nil, fmt.Errorf("failed to build github client: %w", err)
	}

	return client, nil
}
