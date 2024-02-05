package repository

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepoInterface interface {
	SetRefreshToken(refreshToken string, userID uuid.UUID) error
	GetRefreshToken(refreshToken string) (RefreshToken, error)
}

type RefreshTokenFilter struct {
	token string
}

type DatabaseRepo struct {
	Conn *mongo.Client
}

func (database *DatabaseRepo) SetRefreshToken(refreshToken string, userID uuid.UUID) error {
	collection := database.Conn.Database("auth_db").Collection("refresh_tokens")
	_, err := collection.InsertOne(context.Background(), RefreshToken{Token: refreshToken, UserID: userID})
	return err
}

func (database *DatabaseRepo) GetRefreshToken(refreshToken string) (RefreshToken, error) {
	collection := database.Conn.Database("auth_db").Collection("refresh_tokens")
	var token RefreshToken
	err := collection.FindOne(context.Background(), RefreshTokenFilter{token: refreshToken}).Decode(&token)
	return token, err
}
