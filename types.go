package twitch

import (
	"time"
)

type DataResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Total      int         `json:"total"`
}

type StreamsRequest struct {
	After     string   `query:"after"`
	Before    string   `query:"before"`
	First     int      `query:"first"`
	GameID    []string `query:"game_id"`
	Language  string   `query:"language"`
	UserID    []string `query:"user_id"`
	UserLogin []string `query:"user_login"`
}

type UsersFollowsRequest struct {
	After  string `query:"after"`
	First  int    `query:"first"`
	FromID string `query:"from_id"`
	ToID   string `query:"to_id"`
}

type UsersFollows struct {
	FollowedAt time.Time `json:"followed_at"`
	FromID     string    `json:"from_id"`
	FromName   string    `json:"from_name"`
	ToID       string    `json:"to_id"`
	ToName     string    `json:"to_name"`
}

type Stream struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserName     string    `json:"user_name"`
	GameID       string    `json:"game_id"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	Language     string    `json:"language"`
	ThumbnailURL string    `json:"thumbnail_url"`
}

type UsersRequest struct {
	ID    string `query:"id"`
	Login string `query:"login"`
}

type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
}

type Pagination struct {
	Cursor string `json:"cursor"`
}
