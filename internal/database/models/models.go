package models

import "time"

type User struct {
	ID       string `json:"user_id,omitempty"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
}

type Team struct {
	Name    string `json:"team_name"`
	Members []User `json:"members,omitempty"`
}

type PullRequest struct {
	ID                string     `json:"pull_request_id,omitempty"`
	Name              string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"` // OPEN, MERGED
	AssignedReviewers []User     `json:"assigned_reviewers"`
	CreatedAt         time.Time  `json:"created_at"`
	MergedAt          *time.Time `json:"merged_at,omitempty"`
}
