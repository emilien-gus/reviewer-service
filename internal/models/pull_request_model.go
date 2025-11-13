package models

type PullRequest struct {
	ID                string   `db:"pr_id" json:"pr_id"`
	Title             string   `db:"title" json:"title"`
	AuthorID          string   `db:"author_id" json:"author_id"`
	Status            string   `db:"status" json:"status"`
	Reviewers         []string `json:"reviewers"`
	NeedMoreReviewers bool     `db:"need_more_reviewers" json:"need_more_reviewers"`
}
