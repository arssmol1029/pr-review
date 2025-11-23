package models

import "time"

type TeamMember struct {
	UserID   string
	Username string
	IsActive bool
}

type User struct {
	TeamMember
	TeamName string
}

type Team struct {
	Name    string
	Members []TeamMember
}

type PullRequestShort struct {
	ID       string
	Name     string
	AuthorID string
	Status   string
}

type PullRequest struct {
	PullRequestShort
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          *time.Time
}

type UserStats struct {
	UserID        string
	Username      string
	TeamName      string
	OpenReviews   int
	MergedReviews int
	CreatedPRs    int
}

type TeamStats struct {
	TeamName          string
	MemberCount       int
	ActiveMembers     int
	CreatedPRs        int
	AvgReviewersPerPR float64
}

type TotalStats struct {
	TotalTeams        int     `json:"total_teams"`
	TotalUsers        int     `json:"total_users"`
	ActiveUsers       int     `json:"active_users"`
	TotalPRs          int     `json:"total_prs"`
	OpenPRs           int     `json:"open_prs"`
	MergedPRs         int     `json:"merged_prs"`
	AvgReviewersPerPR float64 `json:"avg_reviewers_per_pr"`
}
