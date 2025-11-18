package services

import (
	"context"
	"reviewer-service/internal/models"
	"reviewer-service/internal/repository"
)

type StatsServiceInterface interface {
	GetStats(ctx context.Context) (*models.Stats, error)
}

type StatsService struct {
	repo repository.StatsRepository
}

func NewStatsService(repo repository.StatsRepository) *StatsService {
	return &StatsService{repo: repo}
}

func (s *StatsService) GetStats(ctx context.Context) (*models.Stats, error) {
	users, err := s.repo.GetUserAssignmentsStats(ctx)
	if err != nil {
		return nil, err
	}

	prs, err := s.repo.GetPRAssignmentsStats(ctx)
	if err != nil {
		return nil, err
	}

	return &models.Stats{
		ByUser: users,
		ByPR:   prs,
	}, nil
}
