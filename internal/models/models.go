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
