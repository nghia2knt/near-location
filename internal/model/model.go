package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserLocation struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
	DeletedAt time.Time          `bson:"deleted_at,omitempty" json:"deletedAt"`
	UserId    string             `bson:"user_id" json:"userId"`
	Location  GeoJSON            `bson:"location" json:"location"`
}

type GeoJSON struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type Datapoint struct {
	Longitude float64
	Latitude  float64
}
