package services

import (
	"context"
	"reviewer-service/internal/models"
	"reviewer-service/internal/repository"
)

type PullRequestServiceInterface interface {
	CreatePullRequest(ctx context.Context, prID string, prName string, authorID string) (*models.PullRequest, error)
	SetMergedStatus(ctx context.Context, prid string) (*models.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID string, oldReviewer string) (*models.PullRequest, string, error)
}

type PullRequestService struct {
	PullRequestRepo repository.PullRequestRepositoryInterface
}

func NewPullRequestService(pullRequestRepo repository.PullRequestRepositoryInterface) *PullRequestService {
	return &PullRequestService{PullRequestRepo: pullRequestRepo}
}

func (s *PullRequestService) CreatePullRequest(ctx context.Context, prID string, prName string, authorID string) (*models.PullRequest, error) {
	return s.PullRequestRepo.Create(ctx, prID, prName, authorID)
}

func (s *PullRequestService) SetMergedStatus(ctx context.Context, prid string) (*models.PullRequest, error) {
	return s.PullRequestRepo.SetMerged(ctx, prid)
}

func (s *PullRequestService) ReassignReviewer(ctx context.Context, prID string, oldReviewer string) (*models.PullRequest, string, error) {
	return s.PullRequestRepo.ReassignReviewer(ctx, prID, oldReviewer)
}
