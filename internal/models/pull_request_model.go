package models

import "time"

const (
	StatusOpen   = "OPEN"
	StatusMerged = "MERGED"
)

type PullRequest struct {
	ID        string     `json:"pull_request_id"`
	Name      string     `json:"pull_request_name"`
	AuthorID  string     `json:"author_id"`
	Status    string     `json:"status"`
	Reviewers []string   `json:"assigned_reviewers,omitempty"` // omitempty for ShortPR
	CreatedAt *time.Time `json:"created_at,omitempty"`
	MergedAt  *time.Time `json:"merged_at,omitempty"` // nil if empty
}
