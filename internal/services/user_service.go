package services

import (
	"context"
	"reviewer-service/internal/models"
	"reviewer-service/internal/repository"
)

type UserServiceInterface interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
}

type UserService struct {
	UserRepo repository.UserRepositoryInterface
}

func NewUserService(userRepo repository.UserRepositoryInterface) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	return s.UserRepo.SetIsActive(ctx, userID, isActive)
}
