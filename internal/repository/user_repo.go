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
