package services

import (
	"context"
	"reviewer-service/internal/models"
	"reviewer-service/internal/repository"
)

type TeamServiceInterface interface {
	CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error)
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)
}

type TeamService struct {
	TeamRepo repository.TeamRepositoryInterface
}

func NewTeamService(teamRepo repository.TeamRepositoryInterface) *TeamService {
	return &TeamService{TeamRepo: teamRepo}
}

func (s *TeamService) CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	return s.TeamRepo.Create(ctx, team)
}

func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	return s.TeamRepo.Get(ctx, teamName)
}
