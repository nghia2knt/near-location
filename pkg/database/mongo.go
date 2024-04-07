package database

import (
	"context"
	"near-location/pkg/config"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBConnection(ctx context.Context, databaseName string) (*mongo.Database, error) {
	timeout, err := time.ParseDuration(config.CV.MongoConfig.Timeout)
	if err != nil {
		timeout = time.Second * 30
	}
	maxConIdleTime, err := time.ParseDuration(config.CV.MongoConfig.MaxConnectionIdleTimeOut)
	if err != nil {
		maxConIdleTime = time.Second * 300
	}
	clientOpts := options.Client().
		SetHosts(config.CV.MongoConfig.Address).
		SetConnectTimeout(timeout).
		SetMaxConnIdleTime(maxConIdleTime).
		SetMaxPoolSize(uint64(config.CV.MongoConfig.MaxPoolSize))
	if username, password := config.CV.MongoConfig.Username, config.CV.MongoConfig.Password; username != "" && password != "" {
		auth := options.Credential{
			Username: username,
			Password: password,
		}
		if authSource := config.CV.MongoConfig.AuthSource; authSource != "" {
			auth.AuthSource = authSource
		}
		clientOpts.SetAuth(auth)
	}
	if replica := config.CV.MongoConfig.ReplicaSetName; replica != "" {
		clientOpts.SetReplicaSet(replica)
	}
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	database := client.Database(databaseName)
	log.Info("connected to mongodb")
	return database, nil
}
