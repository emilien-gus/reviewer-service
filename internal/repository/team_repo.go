package repository

import (
	"context"
	"database/sql"
	"errors"
	"reviewer-service/internal/models"
)

var (
	ErrTeamExists   = errors.New(models.CodeTeamExists)
	ErrTeamNotFound = errors.New("team not found")
)

type TeamRepositoryInterface interface {
	Create(ctx context.Context, team *models.Team) (*models.Team, error)
	Get(ctx context.Context, teamName string) (*models.Team, error)
}

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(ctx context.Context, team *models.Team) (*models.Team, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	row := r.db.QueryRowContext(ctx, `
    INSERT INTO teams (name)
    VALUES ($1)
    ON CONFLICT (name) DO NOTHING
    RETURNING name
	`, team.Name)

	err = row.Scan(&team.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTeamExists
		}

		return nil, err
	}

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO users (id, username, team_name, is_active) 
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) DO UPDATE SET 
                username = EXCLUDED.username,
                team_name = EXCLUDED.team_name,
                is_active = EXCLUDED.is_active
    `)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	for _, member := range *team.Members {
		_, err = stmt.ExecContext(ctx,
			member.ID,
			member.Username,
			team.Name,
			member.IsActive,
		)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return team, nil
}

func (r *TeamRepository) Get(ctx context.Context, teamName string) (*models.Team, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT t.name, u.id, u.username, u.is_active
		FROM teams t
		LEFT JOIN users u ON t.name = u.team_name
		WHERE t.name = $1
    `, teamName)

	if err != nil {
		return nil, nil
	}

	defer rows.Close()

	team := &models.Team{}
	var members []models.TeamMember
	teamFound := false

	for rows.Next() {
		teamFound = true
		var member models.TeamMember
		err := rows.Scan(&team.Name, &member.ID, &member.Username, &member.IsActive)
		if err != nil {
			return nil, err
		}

		if member.ID != "" {
			members = append(members, member)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if !teamFound {
		return nil, ErrTeamNotFound
	}

	team.Members = &members
	return team, nil
}
