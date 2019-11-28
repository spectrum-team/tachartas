package models

import "go.mongodb.org/mongo-driver/mongo"

import "context"

type DatabaseConfig struct {
	ConnectionString string
	MongoClient      *mongo.Database
	Ctx              context.Context
}
