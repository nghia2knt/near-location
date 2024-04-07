package service

import (
	"context"
	"near-location/internal/model"
	"near-location/pkg/util"
)

//go:generate mockgen -source userService.go -package service -destination userService_mock.go

type UserService interface {
	FindUserLocationsNearDatapoint(ctx context.Context, datapoint model.Datapoint, maxDistance float64, pageSize int64, pageIdx int64) ([]model.UserLocation, int64, error)
}

type userService struct {
	repo UserRepository
}

type UserRepository interface {
	FindNearUserLocation(ctx context.Context, datapoint model.Datapoint, maxDistance float64, limit int64, skip int64) ([]model.UserLocation, int64, error)
}

func NewUserService(repo UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) FindUserLocationsNearDatapoint(ctx context.Context, datapoint model.Datapoint, maxDistance float64, pageSize int64, pageIdx int64) ([]model.UserLocation, int64, error) {
	skip := pageSize * pageIdx
	userLocations, total, err := s.repo.FindNearUserLocation(ctx, datapoint, maxDistance, pageSize, skip)
	if err != nil {
		return nil, 0, util.ErrInternalServerError(err, "failed to query near user")
	}
	return userLocations, total, nil
}
