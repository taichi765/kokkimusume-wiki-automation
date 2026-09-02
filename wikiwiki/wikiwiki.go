package wikiwiki

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

// ApiEndpointBase is base url to send REST APIs to WikiWiki.
const ApiEndpointBase = "https://api.wikiwiki.jp/kokkimusume"

// AuthRequest is a request for `POST https://api.wikiwiki.jp/<wiki-name>/auth`
type AuthRequest struct {
	Password string `json:"password"`
}

// AuthResponse is a response for `POST https://api.wikiwiki.jp/<wiki-name>/auth`
type AuthResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

// GetPageListResponse is a response for `GET https;//api.wikiwiki.jp/<wiki-name>/pages`
type GetPageListResponse struct {
	Pages []GeneralPageInfo `json:"pages"`
}

// GeneralPageInfo is used in [getPageListResponse].
type GeneralPageInfo struct {
	Name      string `json:"name"`
	Timestamp string `json:"timestamp"`
}

// GetPageResponse is a response for `GET https://api.wikiwiki.jp/<wiki-name>/page/<page-name>`
type GetPageResponse struct {
	Page      string `json:"page"`
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}

// PutPageRequest is a request for `PUT https://api.wikiwiki.jp/<wiki-name>/page/<page-name>`
type PutPageRequest struct {
	Source string `json:"source"`
}

// PutPageResponse is a response for `PUT https://api.wikiwiki.jp/<wiki-name>/page/<page-name>`
type PutPageResponse struct {
	Status string `json:"status"`
}

// Client is HTTP client for wikiwiki's REST API.
type Client struct {
	http *http.Client
	// Token acquired from [GetAuthToken].
	tok string
}

// NewClient creates new client using password.
func NewClient(passwd string) (*Client, error) {
	tok, err := GetAuthToken(passwd)
	if err != nil {
		return nil, err
	}

	return &Client{
		http: &http.Client{},
		tok:  tok,
	}, nil
}

// GetAuthToken gets token from wikiwiki's REST API.
func GetAuthToken(passwd string) (string, error) {
	body, err := json.Marshal(AuthRequest{
		Password: passwd,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal auth request: %w", err)
	}

	buf := bytes.NewBuffer(body)

	res, err := http.Post(ApiEndpointBase+"/auth", "application/json", buf)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("something went wrong with authentication: statusCode = %v", res.StatusCode)
	}

	var resJson AuthResponse
	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if resJson.Status != "ok" {
		return "", fmt.Errorf("something went wrong while getting auth token: status = %s", resJson.Status)
	}

	return resJson.Token, nil
}

// GetPageList gets the list of pages from wikiwiki.
func (c *Client) GetPageList() (GetPageListResponse, error) {
	slog.Debug("getting page list")

	endpoint := ApiEndpointBase + "/pages"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		panic("request must be valid")
	}
	req.Header.Add("Authorizarion", "Bearer "+c.tok)

	res, err := c.http.Do(req)
	if err != nil {
		return GetPageListResponse{}, err
	}
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return GetPageListResponse{}, fmt.Errorf("something went wrong while getting page list: statusCode = %v", res.StatusCode)
	}

	var resJson GetPageListResponse
	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		return GetPageListResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	slog.Debug("successfully got page list")

	return resJson, nil
}

// FetchPageContent fetches page content from WikiWiki's REST API.
func (c *Client) FetchPageContent(page string) (string, error) {
	slog.Debug("fetching page content", slog.String("page", page))

	endpoint := ApiEndpointBase + "/page/" + page
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		panic("request must be valid")
	}
	req.Header.Add("Authorization", "Bearer "+c.tok)

	res, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch page content: %w", err)
	}
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("something went wrong while fetching page content: statusCode = %v", res.StatusCode)
	}

	var resJson GetPageResponse
	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	slog.Debug("successfully fetched page content", slog.String("page", page), slog.String("content", resJson.Source))

	return resJson.Source, nil
}

// UpdatePageContent updates page content using WikiWiki's REST API.
func (c *Client) UpdatePageContent(page string, content string) error {
	slog.Debug("updating page content", slog.String("page", page))

	endpoint := ApiEndpointBase + "/page/" + page
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
	req.Header.Add("Authorization", "Bearer "+c.tok)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update content: %w", err)
	}
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}()

	var resJson PutPageResponse
	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if resJson.Status != "ok" {
		return fmt.Errorf("something went wrong while updating page content: status was %v", resJson.Status)
	}

	slog.Debug("successfully updated page", slog.String("page", page))

	return nil
}
