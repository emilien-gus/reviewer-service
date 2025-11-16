package repository

import (
	"context"
	"database/sql"
	"errors"
	"reviewer-service/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepositoryInterface interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
	GetReviews(ctx context.Context, userID string) (*[]models.PullRequest, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) SetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	var user models.User

	err := r.db.QueryRowContext(ctx, `
        UPDATE users
        SET is_active = $1
        WHERE id = $2
        RETURNING id, username, team_name, is_active
    `, isActive, userID).
		Scan(&user.ID, &user.Username, &user.TeamName, &user.IsActive)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetReviews(ctx context.Context, userID string) (*[]models.PullRequest, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, name, author_id, status
        FROM pull_requests
        WHERE assigned_reviewers @> $1::jsonb
    `, `["`+userID+`"]`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.PullRequest

	for rows.Next() {
		var pr models.PullRequest
		err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &reviews, nil
}
