package main

// Request for `POST https://api.wikiwiki.jp/<wiki-name>/auth`
type AuthRequest struct {
	Password string `json:"password"`
}

// Response for `POST https://api.wikiwiki.jp/<wiki-name>/auth`
type AuthResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

// Response for `GET https://api.wikiwiki.jp/<wiki-name>/page/<page-name>`
type GetPageResponse struct {
	Page      string `json:"page"`
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}

// Request for `PUT https://api.wikiwiki.jp/<wiki-name>/page/<page-name>`
type PutPageRequest struct {
	Source string `json:"source"`
}
