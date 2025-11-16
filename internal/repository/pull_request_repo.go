package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"reviewer-service/internal/models"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

var (
	ErrPullRequestExists   = errors.New(models.CodePRExists)
	ErrPullRequestNotFound = errors.New("pull request not found")
	ErrPRMerged            = errors.New(models.CodePRMerged)
	ErrNotAssigned         = errors.New(models.CodeNotAssigned)
	ErrNoCandidate         = errors.New(models.CodeNoCandidate)
)

type PullRequestRepositoryInterface interface {
	Create(ctx context.Context, prID string, prName string, authorID string) (*models.PullRequest, error)
	SetMerged(ctx context.Context, prId string) (*models.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID string, oldReviewer string) (*models.PullRequest, string, error)
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
		Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &reviewersJSON, &pr.CreatedAt)

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
	log.Print(prId)
	var pr models.PullRequest
	var reviewersJSON []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT * FROM pull_requests
		WHERE id = $1
	`, prId).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &reviewersJSON, &pr.MergedAt, &pr.CreatedAt)

	if err != nil {
		return nil, err
	}

	if pr.ID == "" {
		return nil, ErrPullRequestNotFound
	}

	if pr.Status == models.StatusMerged {
		pr.Reviewers, err = decodeReviewersField(reviewersJSON)
		if err != nil {
			return nil, err
		}
		return &pr, nil
	}

	now := time.Now()
	err = r.db.QueryRowContext(ctx, `
        UPDATE pull_requests 
        SET status = $1, merged_at = $2 
        WHERE id = $3
        RETURNING id, name, author_id, status, assigned_reviewers, merged_at, created_at
    `, models.StatusMerged, now, prId).
		Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &reviewersJSON, &pr.MergedAt, &pr.CreatedAt)

	if err != nil {
		return nil, err
	}

	pr.Reviewers, err = decodeReviewersField(reviewersJSON)
	if err != nil {
		return nil, err
	}

	return &pr, nil
}

func (r *PullRequestRepository) ReassignReviewer(ctx context.Context, prID string, oldReviewer string) (*models.PullRequest, string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback()

	var pr models.PullRequest
	var reviewersJSON []byte

	err = tx.QueryRowContext(ctx, `
		SELECT author_id, status, assigned_reviewers
		FROM pull_requests
		WHERE id = $1
	`, prID).Scan(
		&pr.AuthorID, &pr.Status, &reviewersJSON,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", ErrPullRequestNotFound
		}
		return nil, "", err
	}

	reviewers, err := decodeReviewersField(reviewersJSON)
	if err != nil {
		return nil, "", err
	}

	found := false
	for _, r := range reviewers {
		if r == oldReviewer {
			found = true
			break
		}
	}
	if !found {
		return nil, "", ErrNotAssigned
	}

	var teamName string
	err = tx.QueryRowContext(ctx, `
		SELECT team_name FROM users WHERE id = $1
	`, pr.AuthorID).Scan(&teamName)
	if err != nil {
		return nil, "", err
	}

	qb := squirrel.
		Select("id").
		From("users").
		Where(squirrel.Eq{"team_name": teamName}).
		Where(squirrel.Eq{"is_active": true}).
		Where(squirrel.NotEq{"id": pr.AuthorID}).
		OrderBy("random()").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar)

	for _, rev := range reviewers {
		qb = qb.Where(squirrel.NotEq{"id": rev})
	}

	sqlSelect, args, err := qb.ToSql()
	if err != nil {
		return nil, "", err
	}

	var newReviewer string
	err = tx.QueryRowContext(ctx, sqlSelect, args...).Scan(&newReviewer)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", ErrNoCandidate
		}
		return nil, "", err
	}

	for i, r := range reviewers {
		if r == oldReviewer {
			reviewers[i] = newReviewer
			break
		}
	}

	newJSON, err := json.Marshal(reviewers)
	if err != nil {
		return nil, "", err
	}

	err = tx.QueryRowContext(ctx, `
		UPDATE pull_requests
		SET assigned_reviewers = $1
		WHERE id = $2
		RETURNING id, name, author_id, status, assigned_reviewers, merged_at, created_at
	`, newJSON, prID).
		Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &reviewersJSON, &pr.MergedAt, &pr.CreatedAt)

	if err != nil {
		return nil, "", err
	}

	pr.Reviewers, err = decodeReviewersField(reviewersJSON)
	if err != nil {
		return nil, "", err
	}

	// 6. Коммит
	if err := tx.Commit(); err != nil {
		return nil, "", err
	}

	return &pr, newReviewer, nil
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
