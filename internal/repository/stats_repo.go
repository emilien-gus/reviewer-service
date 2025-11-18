package repository

import (
	"context"
	"database/sql"
)

type StatsRepositoryInterface interface {
	GetUserAssignmentsStats(ctx context.Context) (map[string]int, error)
	GetPRAssignmentsStats(ctx context.Context) (map[string]int, error)
}

type StatsRepository struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) GetUserAssignmentsStats(ctx context.Context) (map[string]int, error) {
	query := `
        SELECT reviewer_id, COUNT(*) AS assignments_count
        FROM (
            SELECT jsonb_array_elements_text(assigned_reviewers) AS reviewer_id
            FROM pull_requests
        ) AS reviewers
        GROUP BY reviewer_id
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)

	for rows.Next() {
		var reviewerID string
		var count int
		if err := rows.Scan(&reviewerID, &count); err != nil {
			return nil, err
		}
		stats[reviewerID] = count
	}

	return stats, rows.Err()
}

func (r *StatsRepository) GetPRAssignmentsStats(ctx context.Context) (map[string]int, error) {
	query := `
        SELECT id, jsonb_array_length(assigned_reviewers) AS reviewers_count
        FROM pull_requests
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)

	for rows.Next() {
		var prID string
		var count int
		if err := rows.Scan(&prID, &count); err != nil {
			return nil, err
		}
		stats[prID] = count
	}

	return stats, rows.Err()
}
