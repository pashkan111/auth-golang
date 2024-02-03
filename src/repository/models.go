package repository

import "github.com/google/uuid"

type RefreshToken struct {
	Token  string
	UserID uuid.UUID
}
