package main

type AuthRequest struct {
	Password string `json:"password"`
}

type AuthResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

type GetPageResponse struct {
	Page      string `json:"page"`
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}
