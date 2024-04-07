package repo

import (
	"context"
	"near-location/internal/model"
	"near-location/pkg/config"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repo struct {
	db *mongo.Database
}

func NewRepo(db *mongo.Database) *repo {
	return &repo{
		db: db,
	}
}

const (
	pointLocationType = "Point"
)

func (r *repo) FindNearUserLocation(ctx context.Context, datapoint model.Datapoint, maxDistance int64, limit int64, skip int64) ([]model.UserLocation, error) {
	nearQuery := bson.M{
		"$geometry": bson.M{
			"type":        pointLocationType,
			"coordinates": []float64{datapoint.Longitude, datapoint.Latitude},
		},
	}
	if maxDistance > 0 {
		nearQuery["$maxDistance"] = maxDistance
	}
	query := bson.M{
		"location": bson.M{
			"$near": nearQuery,
		},
		"deleted_at": bson.M{"$exists": false},
	}
	projection := bson.M{
		"_id":                  1,
		"user_id":              1,
		"location.coordinates": 1,
		"updated_at":           1,
	}
	queryOptions := options.Find().SetLimit(limit).SetSkip(skip).SetProjection(projection)
	cursor, err := r.db.Collection(config.CV.Collection.UserLocation).Find(ctx, query, queryOptions)
	if err != nil {
		log.Errorf("failed to find user location: %s", err)
		return nil, err
	}
	var userLocations []model.UserLocation
	for cursor.Next(ctx) {
		var userLocation model.UserLocation
		if err := cursor.Decode(&userLocation); err != nil {
			log.Errorf("failed to decode user location: %s", err)
			continue
		}
		userLocations = append(userLocations, userLocation)
	}
	return userLocations, nil
}
