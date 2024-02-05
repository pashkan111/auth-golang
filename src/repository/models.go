package repository

import "github.com/google/uuid"

type RefreshToken struct {
	Token  string    `bson:"token"`
	UserID uuid.UUID `bson:"user_id"`
}
