package service

import (
	"context"
	"near-location/internal/model"
	"near-location/pkg/config"
	"near-location/pkg/util"

	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source userService.go -package service -destination userService_mock.go

type UserService interface {
	FindUserLocationsNearDatapoint(ctx context.Context, datapoint model.Datapoint, maxDistance int64, limit int64, skip int64) ([]model.UserLocation, error)
}

type userService struct {
	repo UserRepository
}

type UserRepository interface {
	FindNearUserLocation(ctx context.Context, datapoint model.Datapoint, maxDistance int64, limit int64, skip int64) ([]model.UserLocation, error)
}

func NewUserService(repo UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) FindUserLocationsNearDatapoint(ctx context.Context, datapoint model.Datapoint, maxDistance int64, limit int64, skip int64) ([]model.UserLocation, error) {
	maxLimit := int64(config.CV.MongoConfig.QueryMaxLimit)
	if limit > maxLimit {
		log.Warnf("invalid query limit, %d larger than max limit config %d, auto set to max limit", limit, maxLimit)
		limit = maxLimit
	}
	userLocations, err := s.repo.FindNearUserLocation(ctx, datapoint, maxDistance, limit, skip)
	if err != nil {
		return nil, util.ErrInternalServerError(err, "failed to query near user")
	}
	return userLocations, nil
}
