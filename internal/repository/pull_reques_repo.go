package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"reviewer-service/internal/models"
	"time"

	"github.com/lib/pq"
)

var ErrPullRequestExists = errors.New(models.CodePRExists)
var ErrPullRequestNotFound = errors.New("pull request not found")

type PullRequestRepositoryInterface interface {
	Create(ctx context.Context, prID string, prName string, authorID string) (*models.PullRequest, error)
	SetMerged(ctx context.Context, prId string) (*models.PullRequest, error)
}

type PullRequestRepository struct {
	db *sql.DB
}

func NewPullRequestRepository(db *sql.DB) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}

func (r *PullRequestRepository) Create(ctx context.Context, prID string, prName string, authorID string) (*models.PullRequest, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var teamName sql.NullString
	err = tx.QueryRowContext(ctx, `SELECT team_name FROM users WHERE id = $1`, authorID).Scan(&teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if !teamName.Valid || teamName.String == "" {
		return nil, ErrTeamNotFound
	}

	rows, err := tx.QueryContext(ctx, `
        SELECT id
        FROM users
        WHERE team_name = $1 AND is_active = true AND id <> $2
        ORDER BY random()
        LIMIT 2
    `, teamName.String, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	reviewersJSON, err := json.Marshal(reviewers)
	if err != nil {
		return nil, err
	}
	insertQuery := `
        INSERT INTO pull_requests (id, name, author_id, assigned_reviewers)
        VALUES ($1, $2, $3, $4)
        RETURNING id, name, author_id, status, assigned_reviewers, created_at
    `
	var pr models.PullRequest
	err = tx.QueryRowContext(ctx, insertQuery, prID, prName, authorID, reviewersJSON).
		Scan(&pr.ID, &pr.Name, &pr.AuthorID, &reviewersJSON, &pr.CreatedAt)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return nil, ErrPullRequestExists
		}
		return nil, err
	}

	pr.Reviewers, err = decodeReviewersField(reviewersJSON)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &pr, nil
}

func (r *PullRequestRepository) SetMerged(ctx context.Context, prId string) (*models.PullRequest, error) {
	var pr models.PullRequest
	var reviewersJSON []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT * FROM pull_requests
		WHERE id = $1
	`, prId).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &reviewersJSON, &pr.CreatedAt, &pr.MergedAt)

	if err != nil {
		return nil, err
	}

	if pr.ID == "" {
		return nil, ErrPullRequestNotFound
	}

	if pr.Status == "MERGED" {
		return &pr, nil
	}

	now := time.Now()
	err = r.db.QueryRowContext(ctx, `
        UPDATE pull_requests 
        SET status = $1, merged_at = $2 
        WHERE id = $3
        RETURNING id, name, author_id, status, assigned_reviewers, merged_at, created_at
    `, now, models.StatusMerged, prId).
		Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &reviewersJSON, &pr.CreatedAt, &pr.MergedAt)

	if err != nil {
		return nil, err
	}

	pr.Reviewers, err = decodeReviewersField(reviewersJSON)
	if err != nil {
		return nil, err
	}

	return &pr, nil
}

func decodeReviewersField(reviewersJSON []byte) ([]string, error) {
	var assignedReviewers []string
	if len(reviewersJSON) > 0 {
		if err := json.Unmarshal(reviewersJSON, &assignedReviewers); err != nil {
			return nil, err
		}
	}

	return assignedReviewers, nil
}
