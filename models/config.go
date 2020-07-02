package models

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseConfig struct {
	ConnectionString string
	MongoClient      *mongo.Database
	Ctx              context.Context
}

type AuthInfo struct {
	jwt.StandardClaims
	Email    string
	FullName string
}
